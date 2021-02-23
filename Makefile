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
	$(GOBUILD) -o ./dist/$(BINARY_NAME)
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

dynamodb-admin:
	DYNAMO_ENDPOINT=http://localhost:8000 dynamodb-admin
start-dynamodb:
	docker run \
	--name dynamodb \
	-d \
	-p 8000:8000 amazon/dynamodb-local
list-tables:
	aws dynamodb list-tables \
	--endpoint-url "http://localhost:8000"
create-table:
	aws dynamodb create-table \
	--endpoint-url http://localhost:8000 \
	--table-name 'Gemini' \
	--attribute-definitions \
	AttributeName=PK,AttributeType=S \
	AttributeName=SK,AttributeType=S \
	--key-schema \
	AttributeName=PK,KeyType=HASH \
	AttributeName=SK,KeyType=RANGE \
	--provisioned-throughput \
	ReadCapacityUnits=1,WriteCapacityUnits=1 \
	> /dev/null
delete-table:
	aws dynamodb delete-table \
	--table-name "Gemini" \
	--endpoint-url "http://localhost:8000" > /dev/null
