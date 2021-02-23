package main

import (
	"math/rand"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/segmentio/ksuid"
)

// RandomInt generate a random integer between min, max
func RandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

// GenerateID generate a KSUID with prefix_
func GenerateID(prefix string) string {
	return prefix + "_" + ksuid.New().String()
}

// RandomFloat generate a random float between 0, max
func RandomFloat(max float32) float32 {
	return rand.Float32() * max
}

// RandomElement choose a random element from list
func RandomElement(list []string) string {
	return list[RandomInt(0, len(list))]
}

// DynamoClient get a dynamodb client
func DynamoClient() *dynamodb.DynamoDB {
	endpointURL := "http://docker.for.mac.localhost:8000"

	// FIXME base configs on env
	return dynamodb.New(session.New(), &aws.Config{
		Region:   aws.String("us-west-1"),
		Endpoint: aws.String(endpointURL),
	})
}
