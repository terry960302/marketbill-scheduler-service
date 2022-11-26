package models

import "gorm.io/gorm"

type FlowerBatchProcessLogs struct {
	gorm.Model
	NewFlowerCount      int    `json:"new_flower_count"`          // 새로 추가된 꽃
	NewFlowerTypeCount  int    `json:"new_flower_type_count"`     // 새로 추가된 꽃 품목
	AffectedFlowerCount int    `json:"affected_flower_count"`     // 수정된 꽃(경매 날짜 업데이트)
	Status              string `json:"status"`                    // 성공, 실패 여부
	ErrLogs             string `gorm:"type:text" json:"err_logs"` // 에러인 경우의 로그
}
