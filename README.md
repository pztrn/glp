# GLP

GLP stands for Go Licensing Processor. It is an utility that parses application's dependencies, gets their licenses and copyright information and writing gathered data into report file.

## Supported languages

* Go (dep and modules)

## Supported report file formats

* CSV

## Supported VCS and sites

It was tested for git repositories with various sites (github.com, gitlab.com, self-hosted Gitea, even [my giredore](https://sources.dev.pztrn.name/pztrn/giredore)). It will work with any hosting that supports ``?go-get=1`` URL parameter and which outputs go-import and go-source meta lines.

But there are some caveats appeared:

* Github most of times will not add ``go-source`` meta line in HTML's ``<head>`` tag. There are a workaround for that [here](https://sources.dev.pztrn.name/pztrn/glp/src/branch/master/structs/vcsdata.go).

## Installation

It is enough to issue:

```bash
go get -u go.dev.pztrn.name/glp/cmd/glp
```

## Usage

See `glp -h` for a list of possible options.

### Example usage

```bash
glp  -config ./.glp.yaml -pkgs /home/pztrn/projects/go/src/go.dev.pztrn.name/discordrone,/home/pztrn/projects/go/src/go.dev.pztrn.name/opensaps -outfile /home/pztrn/deps-test.csv
```

## Configuration

For now you can configure only debug output for logging. See ToDo below.

## ToDo

* Ability to overwrite all things about dependency, like copyrights, license URL and so on via configuration file.
* Ability to use it as library.
* Ability to use it in CI with alerts about bad licenses.
* Ability to use it for projects written in other languages than Go (javascript, python,  java, and so on).
* More outputters - PDF, xlsx and so on.
* (Maybe) Use ``go list`` output for gathering dependencies.
