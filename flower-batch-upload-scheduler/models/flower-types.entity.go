package models

import "gorm.io/gorm"

type FlowerTypes struct {
	gorm.Model
	Name string `json:"name"`
	// Flowers []Flowers `gorm:"foreignKey:FlowerTypeID"`
}
