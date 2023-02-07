package datastore

import (
	"fmt"
	"os"
	"public-flower-upload-scheduler/models"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresql() (*gorm.DB, error) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	DSN := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PW") +
		" port=" + os.Getenv("DB_PORT") +
		" database=" + fmt.Sprint(os.Getenv("DB_NAME")) +
		" sslmode=disable" +
		" TimeZone=Asia/Seoul"

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	db.AutoMigrate(&models.PublicBiddingFlowers{}, &models.FlowerBatchUploadLogs{}, &models.FlowerBatchProcessLogs{})

	logger.Info("completed")
	return db, nil
}
