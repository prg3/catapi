# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=catapi
BINARY_UNIX=$(BINARY_NAME)_unix

all: deps build docker
deps:
	$(GOGET) -u github.com/onsi/ginkgo
	$(GOGET) -u github.com/go-redis/redis
build: 
	$(GOBUILD) -a -installsuffix cgo -o catapi .
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
docker:
	docker build -t catapi .
