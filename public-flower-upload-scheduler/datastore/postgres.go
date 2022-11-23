package datastore

import (
	"log"
	"os"
	"public-flower-upload-scheduler/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresql() *gorm.DB {
	DSN := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PW") +
		" port=" + os.Getenv("DB_PORT") +
		" database=" + os.Getenv("DB_NAME") +
		" sslmode=disable" +
		" TimeZone=Asia/Seoul"

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	db.AutoMigrate(&models.PublicBiddingFlowers{}, &models.FlowerUploadLogs{})

	return db
}
