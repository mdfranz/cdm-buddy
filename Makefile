# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=cdmbuddy
BINARY_UNIX=$(BINARY_NAME)_unix

# Directory structure
CMD_DIR=./cmd/cdmbuddy

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(CMD_DIR)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f *.csv *.json *.xlsx

run:
	$(GOCMD) run $(CMD_DIR)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../$(BINARY_UNIX) $(CMD_DIR)

.PHONY: all build test clean run build-linux
