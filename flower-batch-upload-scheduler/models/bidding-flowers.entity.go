package models

import (
	"time"

	"gorm.io/gorm"
)

type BiddingFlowers struct {
	gorm.Model
	FlowerID    uint      `json:"flower_id"`
	BiddingDate time.Time `json:"bidding_date"`
}
