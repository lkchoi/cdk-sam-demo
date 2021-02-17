# Go parameters
GOOS=linux
GOARCH=amd64
GOENV=GOOS=$(GOOS) GOARCH=$(GOARCH)
GOCMD=go
GOBUILD=$(GOENV) $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
CORE_BINARY_NAME=aws-lambda-go-api-proxy-core
GIN_BINARY_NAME=aws-lambda-go-api-proxy-gin
SAMPLE_BINARY_NAME=main

all: clean test build package
build:
	$(GOBUILD) ./...
	cd app && $(GOBUILD) -o $(SAMPLE_BINARY_NAME)
package:
	mkdir -p dist && mv app/$(SAMPLE_BINARY_NAME) dist
test:
	$(GOTEST) -v ./...
synth:
	cd cdk && cdk synth --no-staging > template.yaml
invoke: synth
	cd cdk && sam local invoke
start:
	cd cdk && sam local start-api
clean:
	rm -f dist/$(SAMPLE_BINARY_NAME)
