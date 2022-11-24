package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"public-flower-upload-scheduler/datastore"
	"public-flower-upload-scheduler/models"
	"public-flower-upload-scheduler/services"

	"github.com/aws/aws-lambda-go/events"
	"gorm.io/gorm"
)

func HandleUpload(r *models.LambdaResponse, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := datastore.NewPostgresql()
	if err != nil {
		e := fmt.Sprintf("[NewPostgresql] %s", err.Error())
		return r.Error(http.StatusInternalServerError, e)

	}
	var errLog models.FlowerUploadLogs

	defer func() {
		if err := recover(); err != nil {
			errLog = models.FlowerUploadLogs{
				Success: -1,
				Failure: -1,
				Total:   -1,
				ErrLogs: fmt.Sprint(err),
			}
		}
	}()

	log, err := UploadFlowers(db, request)
	if err != nil {
		errLog = models.FlowerUploadLogs{
			Success: -1,
			Failure: -1,
			Total:   -1,
			ErrLogs: err.Error(),
		}
		bytes, _ := json.Marshal(errLog)
		db.Create(errLog)
		return r.Error(http.StatusInternalServerError, string(bytes))
	}

	return r.Json(http.StatusOK, log)
}

func UploadFlowers(db *gorm.DB, request events.APIGatewayProxyRequest) (*models.FlowerUploadLogs, error) {

	flowers, err := services.FetchFlowerItems()
	if err != nil {
		return nil, err
	}

	log, err := services.UploadRawFlowerData(db, flowers)
	if err != nil {
		return nil, err
	}

	tx := db.Create(log)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return log, nil
}
