package models

import (
	"time"

	"gorm.io/gorm"
)

type PublicBiddingFlower struct {
	gorm.Model
	FlowerType string    `json:"flower_type"`
	FlowerName string    `json:"flower_name"`
	Grade      string    `json:"grade"`
	Quantity   int       `json:"quantity"`
	MaxPrice   string    `json:"max_price"`
	MinPrice   string    `json:"min_price"`
	AvgPrice   string    `json:"avg_price"`
	TotalPrice string    `json:"total_price"`
	BidDate    time.Time `json:"bid_date"`
}
