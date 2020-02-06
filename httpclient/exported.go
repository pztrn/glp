package httpclient

import (
	// stdlib
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	// local
	"go.dev.pztrn.name/glp/configuration"
)

const (
	defaultTimeoutInSeconds = 20
	perDomainRequestsLimit  = 5
)

var (
	httpClient = &http.Client{
		Timeout: time.Second * defaultTimeoutInSeconds,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   time.Second * defaultTimeoutInSeconds,
				DualStack: true,
			}).DialContext,
			ExpectContinueTimeout: time.Second * 5,
			Proxy:                 http.ProxyFromEnvironment,
			ResponseHeaderTimeout: time.Second * defaultTimeoutInSeconds,
			TLSHandshakeTimeout:   time.Second * 5,
		},
	}

	perDomainRequests      map[string]int
	perDomainRequestsMutex sync.Mutex
)

// Initialize initializes package.
func Initialize() {
	log.Println("Initializing HTTP client...")

	perDomainRequests = make(map[string]int)
}

// GET executes GET request and returns body.
func GET(request *http.Request) []byte {
	for {
		perDomainRequestsMutex.Lock()
		currentlyRunning, found := perDomainRequests[request.URL.Host]
		perDomainRequestsMutex.Unlock()

		if !found {
			break
		}

		if currentlyRunning >= perDomainRequestsLimit {
			time.Sleep(time.Second * 1)
			continue
		}

		break
	}

	perDomainRequestsMutex.Lock()

	_, found := perDomainRequests[request.URL.Host]
	if !found {
		perDomainRequests[request.URL.Host] = 1
	} else {
		perDomainRequests[request.URL.Host]++
	}

	perDomainRequestsMutex.Unlock()

	defer func() {
		perDomainRequestsMutex.Lock()

		perDomainRequests[request.URL.Host]--

		perDomainRequestsMutex.Unlock()
	}()

	if configuration.Cfg.Log.Debug {
		log.Println("Executing request:", request.URL.String())
	}

	var (
		requestsCount = 0
		response      *http.Response
	)

	for {
		if requestsCount == 3 {
			log.Printf("Failed to execute request %s: tried 3 times and got errors. Skipping.", request.URL.String())
			return nil
		}

		var err error

		response, err = httpClient.Do(request)
		if err != nil {
			log.Printf("Failed to execute request %s: %s\n", request.URL.String(), err.Error())
			requestsCount++
			time.Sleep(time.Second * 1)
			continue
		}

		break
	}

	respBody, err1 := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err1 != nil {
		log.Printf("Failed to read response body %s: %s\n", request.URL.String(), err1.Error())
		return nil
	}

	return respBody
}
