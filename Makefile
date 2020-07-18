# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get
BINARY_DIR=./bin
BINARY_NAME=$(BINARY_DIR)/sqs_wrapper
BINARY_UNIX=$(BINARY_NAME)_unix


all: test build ## Build and test the binary

build: ## Build the binary
	$(GOBUILD) -o $(BINARY_NAME) -v

test: ## Test all the test files recursively
	$(GOTEST) -v ./tests/... -coverpkg=./...

test-cover: ## Test and generate the coverage report
	$(GOTEST)  -coverprofile=coverage.out ./tests/... -coverpkg=./... && $(GOTOOL) cover -func=coverage.out

test-show-cover:
	$(GOTOOL) cover -html=coverage.out

clean: ## Clean the binaries
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

.PHONY: run
run: ## Run the binary
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

.PHONY: build-linux
build-linux: ## Cross compilation
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

.PHONY: --help
--help: ##
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: help
help: --help
