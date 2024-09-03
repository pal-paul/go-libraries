SERVICE		?= $(shell basename `go list`)
VERSION		?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(PWD)/.version 2> /dev/null || echo v0)
PACKAGE		?= $(shell go list)
PACKAGES	?= $(shell go list ./...)
FILES		?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: help clean fmt lint vet test build all

default: help

help:   ## show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

all:    ## clean, format, build and unit test
	make clean-all
	make build
	make test

install:    ## build and install go application executable
	go install -v ./...
	go install go.uber.org/mock/mockgen@latest
	go install github.com/segmentio/golines@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

env:    ## Print useful environment variables to stdout
	echo $(CURDIR)
	echo $(SERVICE)
	echo $(PACKAGE)
	echo $(VERSION)

clean:  ## go clean
	go clean

clean-all:  ## remove all generated artifacts and clean all build artifacts
	go clean -i ./...

tools:  ## fetch and install all required tools

vet:    ## run go vet on the source files
	go vet ./...

doc:    ## generate godocs and start a local documentation webserver on port 8085

update-dependencies:    ## update golang dependencies
	dep ensure

generate-mocks:     ## generate mock code
	go generate ./...

build: generate-mocks ## generate all mocks and build the go code

deploy: install build

test: ## run tests
	go test -v ./...

tidy:
	go get -u ./...
	go mod tidy