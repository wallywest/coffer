current_dir := $(patsubst %/,%, $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
REPO_PATH := gitlab.vailsys.com/jerny/coffer
VERSION=$(shell cat $(current_dir)/version/VERSION)
REV := $(shell git rev-parse --short HEAD 2> /dev/null  || echo 'unknown')
BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown')
BUILD_DATE := $(shell date +%Y%m%d-%H:%M:%S)
BUILDFLAGS := -ldflags \
			 " -X $(REPO_PATH)/version.Version=$(VERSION)\
			   -X $(REPO_PATH)/version.Revision=$(REV)\
			   -X $(REPO_PATH)/version.Branch=$(BRANCH)\
			   -X $(REPO_PATH)/version.BuildDate=$(BUILD_DATE)"

BUILD_PATH := gitlab.vailsys.com/jerny/coffer/cmd/coffer
DOCKER_IP := 192.168.99.100
EXPLORATORY_SERVER := spv07vcs16.vail

setup:
	@go get -u "github.com/tools/godep"
	@go get -u "github.com/golang/lint/golint"

lint:
	@echo "==== go lint ==="
	@golint ./...

vet:
	@echo "=== go vet ==="
	@go vet ./...

fmt:
	@echo "=== go fmt ==="
	@go fmt ./...

test: fmt
	@echo "=== go test ==="
	@godep go test ./... -cover

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
