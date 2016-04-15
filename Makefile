current_dir := $(patsubst %/,%, $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
REPO_PATH=${BUILD_PATH}/${PACKAGE_NAME}
BUILD_PATH=${REPO_PATH}/cmd/coffer
REPO_PATH = gitlab.vailsys.com/jerny/
PACKAGE_NAME = coffer
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

setup:
	@go get -u "github.com/tools/godep"
	@go get -u "github.com/golang/lint/golint"

list:
	@echo $(GOFILES_NOVENDOR)

lint:
	@echo "==== go lint ==="
	@golint $(GOFILES_NOVENDOR)

vet:
	@echo "=== go vet ==="
	@go vet $(GOFILES_NOVENDOR)

fmt:
	@echo "=== go fmt ==="
	@gofmt -l -w ${GOFILES_NOVENDOR}

test: fmt
	@echo "=== go test ==="
	@godep go test $(GOFILES_NOVENDOR) -cover -short

update-deps:
	@echo "=== godep update ==="
	@godep save -r ./...

install: test	
	@echo "=== go install ==="
	@godep go install $(BUILDFLAGS) $(BUILD_PATH)

build: 
	@echo "=== go build ==="
	@mkdir -p bin/
	@godep go build -o bin/coffer $(BUILDFLAGS) $(BUILD_PATH)

docker-build:
	@docker build --rm=true --no-cache=true --file=Dockerfile.dev -t coffer:dev .

run-dev: build
	./bin/coffer \
	--loglevel DEBUG \
	--env test \
	--skip-registration \
	--consul consul://$(DOCKER_IP):8500 \
	--vcsdbm $(DOCKER_IP):8255 \
	--vcsserver $(DOCKER_IP):8265 \
	--discovery static \
	start

run-integration: build
	./bin/coffer \
	--loglevel DEBUG \
	--env test \
	--consul consul://$(EXPLORATORY_SERVER):8500 \
	--vcsdbm $(EXPLORATORY_SERVER):8255 \
	--vcsserver $(EXPLORATORY_SERVER):8265 \
	--discovery static \
	start

.PHONY: setup lint vet install test build update-deps
