package models

import "gorm.io/gorm"

type FlowerUploadLogs struct {
	gorm.Model
	Success int    `json:"success"`
	Failure int    `json:"failure"`
	Total   int    `json:"total"`
	ErrLogs string `gorm:"type:text"`
}
