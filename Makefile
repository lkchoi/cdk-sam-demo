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
BINARY_NAME=main

all: clean test build
build:
	$(GOBUILD) ./...
	mkdir -p dist
	cd app && $(GOBUILD) -o ../dist/$(BINARY_NAME)
test:
	$(GOTEST) -v ./...
synth:
	cd cdk && cdk synth --no-staging > template.yaml
invoke: synth
	cd cdk && sam local invoke
start:
	cd cdk && sam local start-api
diff:
	cd cdk && cdk diff
deploy: build docs synth
	cd cdk && cdk deploy
clean:
	rm -f dist/$(BINARY_NAME)
