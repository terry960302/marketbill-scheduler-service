package models

import "gorm.io/gorm"

type FlowerBatchUploadLogs struct {
	gorm.Model
	Success int    `json:"success"`                   // 업로드 성공개수
	Failure int    `json:"failure"`                   // 업로드 실패개수
	Total   int    `json:"total"`                     // 전체 API반환값 개수
	ErrLogs string `gorm:"type:text" json:"err_logs"` // 에러 로그
}
