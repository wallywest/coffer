current_dir := $(patsubst %/,%, $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
BUILD_PATH=gitlab.vailsys.com/jerny/coffer
REPO_PATH=${BUILD_PATH}
GOPACKAGES?=${BUILD_PATH}/${PACKAGE_NAME}/...
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
VERSION=$(shell cat $(current_dir)/version/VERSION)
REV := $(shell git rev-parse --short HEAD 2> /dev/null  || echo 'unknown')
BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown')
BUILD_DATE := $(shell date +%Y%m%d-%H:%M:%S)
BUILDFLAGS := -ldflags \
			 " -X $(REPO_PATH)/version.Version=$(VERSION)\
			   -X $(REPO_PATH)/version.Revision=$(REV)\
			   -X $(REPO_PATH)/version.Branch=$(BRANCH)\
			   -X $(REPO_PATH)/version.BuildDate=$(BUILD_DATE)"

tools:
	@go get github.com/onsi/ginkgo/ginkgo
	@go get github.com/onsi/gomega 
	@go get -u github.com/kardianos/govendor

test:
	@ginkgo -r 

test-ci:
	@ginkgo -r -noColor -succinct

update-deps:
	@echo "=== govendor update ==="
	@govendor update +vendor

build:
	@echo "=== go build ==="
	@mkdir -p bin/
	@govendor build -o bin/coffer $(BUILDFLAGS) $(BUILD_PATH)

.PHONY: setup lint vet install test build update-deps
