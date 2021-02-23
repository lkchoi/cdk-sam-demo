package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

// Handler is the main entry point for Lambda. Receives a proxy request and
// returns a proxy response
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		log.Printf("Gin cold start")
		r := gin.Default()

		r.GET("/api/v1/spaces", getSpaces)
		r.GET("/api/v1/spaces/:id", getSpace)
		r.POST("/api/v1/spaces", createSpace)
		r.PUT("/api/v1/spaces/:id", updateSpace)

		r.GET("/api/v1/clients", getClients)
		r.GET("/api/v1/clients/:id", getClient)
		r.POST("/api/v1/clients", createClient)
		r.PUT("/api/v1/clients/:id", updateClient)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
