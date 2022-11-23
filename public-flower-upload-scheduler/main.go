package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"public-flower-upload-scheduler/handlers"
	"public-flower-upload-scheduler/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	profile := os.Getenv("PROFILE")
	log.Print("PROFILE : ", profile)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	r := models.NewLambdaResponse()
	switch request.HTTPMethod {
	case "GET":
		return handlers.HealthCheck(request)
	case "POST":
		return handlers.HandleUpload(request)
	default:
		return r.Error(http.StatusBadRequest, "Wrong http method")
	}
}

func main() {
	lambda.Start(HandleRequest)
}
