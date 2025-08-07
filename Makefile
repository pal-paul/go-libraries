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

deps:
	go mod tidy
	go mod download

update-deps:
	go get -u all
	go mod tidy

tidy:
	go get -u ./...
	go mod tidy

vet:    ## run go vet on the source files
	go vet ./...

doc:    ## generate godocs and start a local documentation webserver on port 8085

lint:
	golangci-lint run --verbose

update-dependencies:    ## update golang dependencies
	dep ensure

generate-mocks:     ## generate mock code
	go generate ./...

build: generate-mocks ## generate all mocks and build the go code

deploy: install build

test: ## run tests
	go test -v ./...

# HTTP Server targets
test-http-server: ## run HTTP server tests
	go test -v ./pkg/http-server

benchmark-http-server: ## run HTTP server benchmarks
	go test -bench=. -benchmem ./pkg/http-server

run-http-server-example: ## run the HTTP server example
	go run examples/http-server/main.go

run-benchmark: ## run the performance benchmark tool (requires server to be running)
	cd examples/benchmark && go run main.go

build-examples: ## build all examples
	cd examples/http-server && go build -o http-server main.go
	cd examples/benchmark && go build -o benchmark main.go
	cd examples/comparison && go build servers.go && go build compare.go

# Performance comparison targets
run-server-comparison: ## run comprehensive HTTP server performance comparison
	cd examples/comparison && ./run_comparison.sh

run-server-comparison-keep: ## run comparison and keep servers running for manual testing
	cd examples/comparison && ./run_comparison.sh --keep-running

start-fasthttp-server: ## start only our FastHTTPServer on port 8080
	cd examples/comparison && go run servers.go fasthttp-server

start-chi-server: ## start only Go-Chi server on port 8081
	cd examples/comparison && go run servers.go chi

start-std-server: ## start only Standard HTTP server on port 8082
	cd examples/comparison && go run servers.go std

start-fasthttp-lib-server: ## start only Fasthttp library server on port 8083
	cd examples/comparison && go run servers.go fasthttp