# GLP

GLP stands for Go Licensing Processor. It is an utility that parses application's dependencies, gets their licenses and copyright information and writing gathered data into report file.

## Supported languages

* Go (dep and modules)

## Supported report file formats

* CSV

## Supported VCS and sites

None yet. It executes HTTP request with ``?go-get=1`` parameter to get go-import and go-source data.

## Installation

It is enough to issue:

```bash
go get -u go.dev.pztrn.name/glp/cmd/glp
```

## Configuration

*None yet.*
