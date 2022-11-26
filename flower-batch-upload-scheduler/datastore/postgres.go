package datastore

import (
	"fmt"
	"log"
	"os"
	"public-flower-upload-scheduler/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresql() (*gorm.DB, error) {
	DSN := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PW") +
		" port=" + os.Getenv("DB_PORT") +
		" database=" + fmt.Sprint(os.Getenv("DB_NAME")) +
		" sslmode=disable" +
		" TimeZone=Asia/Seoul"

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		log.Print(err)
		return nil, err
	}

	db.AutoMigrate(&models.PublicBiddingFlowers{}, &models.FlowerBatchUploadLogs{})

	return db, nil
}
