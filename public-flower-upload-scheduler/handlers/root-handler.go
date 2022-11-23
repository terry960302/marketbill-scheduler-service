package handlers

import (
	"fmt"
	"net/http"
	"os"
	"public-flower-upload-scheduler/models"

	"github.com/aws/aws-lambda-go/events"
)

func HealthCheck(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r := models.NewLambdaResponse()
	profile := os.Getenv("PROFILE")
	msg := fmt.Sprintf("[%s] Public Flower Upload Scheduler Service is running...", profile)
	return r.Json(http.StatusOK, msg)
}
