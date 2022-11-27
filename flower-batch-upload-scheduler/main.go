package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"public-flower-upload-scheduler/datastore"
	"public-flower-upload-scheduler/handlers"
	"public-flower-upload-scheduler/models"

	"github.com/aws/aws-lambda-go/events"
)

func init() {
	profile := os.Getenv("PROFILE")
	log.Print("PROFILE : ", profile)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r := models.NewLambdaResponse()
	log.Printf("[LOG] method : %s", request.HTTPMethod)
	switch request.HTTPMethod {
	case "GET":
		return handlers.HealthCheck(r, request)
	case "POST":
		return handlers.HandleUpload(r, request)
	default:
		return r.Error(http.StatusBadRequest, "Wrong http method")
	}
}

// func main() {
// 	lambda.Start(HandleRequest)
// }

// test
func main() {
	os.Setenv("PROFILE", "dev")
	os.Setenv("DB_USER", "marketbill")
	os.Setenv("DB_PW", "marketbill1234!")
	os.Setenv("DB_NET", "tcp")
	os.Setenv("DB_HOST", "marketbill-db.ciegftzvpg1l.ap-northeast-2.rds.amazonaws.com")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "dev-db")
	os.Setenv("PORT", "8080")
	os.Setenv("API_KEY", "4DC6A10B4F5D43D5977F364FC0DFE81C")

	db, _ := datastore.NewPostgresql()
	db.AutoMigrate(&models.FlowerBatchUploadLogs{})

	// flowers := []models.Flowers{}
	// db.Table("flowers").Joins("FlowerTypes").Find(&flowers)

	// utils.PrettyPrint(flowers)

}
