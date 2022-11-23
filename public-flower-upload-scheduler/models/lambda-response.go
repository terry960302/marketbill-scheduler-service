package models

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type LambdaResponse struct {
}

func NewLambdaResponse() *LambdaResponse {
	return &LambdaResponse{}
}

func (r *LambdaResponse) Error(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept"}
	return events.APIGatewayProxyResponse{
		StatusCode:        statusCode,
		Headers:           headers,
		MultiValueHeaders: nil,
		Body:              message,
		IsBase64Encoded:   false}, nil
}

func (r *LambdaResponse) Json(statusCode int, body any) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept"}

	bytes, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:        http.StatusInternalServerError,
			Headers:           headers,
			MultiValueHeaders: nil,
			Body:              err.Error(),
			IsBase64Encoded:   false}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode:        statusCode,
		Headers:           headers,
		MultiValueHeaders: nil,
		Body:              string(bytes),
		IsBase64Encoded:   false}, nil
}
