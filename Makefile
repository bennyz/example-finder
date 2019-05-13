# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

PREFIX=.
ARTIFACT_DIR ?= .

VERSION?=$(shell git describe --tags --always --match "v[0-9]*" | awk -F'-' '{print substr($$1,2) }')
RELEASE?=$(shell git describe --tags --always --match "v[0-9]*" | awk -F'-' '{if ($$2 != "") {print $$2 "." $$3} else {print 1}}')
VERSION_RELEASE=$(VERSION)$(if $(RELEASE),-$(RELEASE))

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

COMMON_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COMMON_GO_BUILD_FLAGS=
# COMMON_GO_BUILD_FLAGS=-ldflags '-extldflags "-static"'

all: clean build

build: fmt test vet
	$(COMMON_ENV) $(GOBUILD) $(COMMON_GO_BUILD_FLAGS)
    	
.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

test:
	$(GOTEST) -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o /tmp/coverage.html

clean:
	$(GOCLEAN)
	git clean -dfx -e .idea*

.PHONY: all test build
