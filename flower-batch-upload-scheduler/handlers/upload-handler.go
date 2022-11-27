package handlers

import (
	"fmt"
	"net/http"
	"public-flower-upload-scheduler/datastore"
	"public-flower-upload-scheduler/models"
	"public-flower-upload-scheduler/services"

	"github.com/aws/aws-lambda-go/events"
)

func HandleUpload(r *models.LambdaResponse, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := datastore.NewPostgresql()
	if err != nil {
		e := fmt.Sprintf("[NewPostgresql] %s", err.Error())
		return r.Error(http.StatusInternalServerError, e)
	}

	rawFlowers, err := services.FetchRawFlowerItems()
	if err != nil {
		return r.Error(http.StatusInternalServerError, err.Error())
	}

	publicFlowers, uploadLog, err := services.UploadFlowers(db, rawFlowers)
	if err != nil {
		return r.Error(http.StatusInternalServerError, err.Error())
	}

	processLog, err := services.ProcessFlowerRawData(db, *publicFlowers)
	if err != nil {
		return r.Error(http.StatusInternalServerError, err.Error())
	}

	return r.Json(http.StatusOK, map[string]interface{}{
		"upload_log":  uploadLog,
		"process_log": processLog,
	})
}
