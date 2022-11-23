package models

import (
	"time"

	"gorm.io/gorm"
)

type PublicBiddingFlowers struct {
	gorm.Model
	FlowerType string    `json:"flower_type"`
	FlowerName string    `json:"flower_name"`
	Grade      string    `json:"grade"`
	Quantity   int       `json:"quantity"`
	MaxPrice   int       `json:"max_price"`
	MinPrice   int       `json:"min_price"`
	AvgPrice   int       `json:"avg_price"`
	TotalPrice int       `json:"total_price"`
	BidDate    time.Time `json:"bid_date"`
}
