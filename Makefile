# Go parameters
GOCMD=go
GOPATH=$(CURDIR)
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=catapi
BINARY_UNIX=$(BINARY_NAME)_unix

all: deps  build build-linux docker
deps:
	$(GOGET) -u github.com/onsi/ginkgo
	$(GOGET) -u github.com/go-redis/redis
build: 
	$(GOINSTALL) github.com/prg3/catapi
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOINSTALL) github.com/prg3/catapi

docker:
	docker build -t catapi .
