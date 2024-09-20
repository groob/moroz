all: build

.PHONY: build

ifndef ($(GOPATH))
	GOPATH = $(HOME)/go
endif

export GO111MODULE=on
PATH := $(GOPATH)/bin:$(PATH)
VERSION = $(shell git describe --tags --always --dirty)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse HEAD)
REVSHORT = $(shell git rev-parse --short HEAD)
USER = $(shell whoami)
GOVERSION = $(shell go version | awk '{print $$3}')
NOW	= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SHELL = /bin/bash
DOCKER_IMAGE_NAME = groob/moroz
DOCKER_IMAGE_TAG = $(shell echo ${VERSION} | sed 's/^v//')

ifneq ($(OS), Windows_NT)
	CURRENT_PLATFORM = linux
	ifeq ($(shell uname), Darwin)
		SHELL := /bin/bash
		CURRENT_PLATFORM = darwin
	endif
else
	CURRENT_PLATFORM = windows
endif

BUILD_VERSION = "\
	-X github.com/kolide/kit/version.appName=${APP_NAME} \
	-X github.com/kolide/kit/version.version=${VERSION} \
	-X github.com/kolide/kit/version.branch=${BRANCH} \
	-X github.com/kolide/kit/version.buildUser=${USER} \
	-X github.com/kolide/kit/version.buildDate=${NOW} \
	-X github.com/kolide/kit/version.revision=${REVISION} \
	-X github.com/kolide/kit/version.goVersion=${GOVERSION}"


gomodcheck: 
	@go help mod > /dev/null || (@echo micromdm requires Go version 1.23 or higher && exit 1)

deps: gomodcheck
	@go mod download

test:
	go test -cover -race -v $(shell go list ./... | grep -v /vendor/)

build: moroz

clean:
	rm -rf build/
	rm -f *.zip

.pre-build:
	mkdir -p build/darwin
	mkdir -p build/linux

.pre-moroz:
	$(eval APP_NAME = moroz)

moroz: .pre-build .pre-moroz
	go build -o build/$(CURRENT_PLATFORM)/moroz -ldflags ${BUILD_VERSION} ./cmd/moroz

xp-moroz: .pre-build .pre-moroz
	GOOS=darwin go build -o build/darwin/moroz -ldflags ${BUILD_VERSION} ./cmd/moroz
	GOOS=linux CGO_ENABLED=0 go build -o build/linux/moroz  -ldflags ${BUILD_VERSION} ./cmd/moroz

install: .pre-moroz
	go install -ldflags ${BUILD_VERSION} ./cmd/moroz

release-zip: xp-moroz
	zip -r moroz_${VERSION}.zip build/

docker-build:
	GOOS=linux CGO_ENABLED=0 go build -o build/linux/moroz  -ldflags ${BUILD_VERSION} ./cmd/moroz
	docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} .

docker-tag: docker-build
	docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ${DOCKER_IMAGE_NAME}:latest
