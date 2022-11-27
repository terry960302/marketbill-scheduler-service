package models

import "gorm.io/gorm"

type Flowers struct {
	gorm.Model
	Name         string `json:"name"`
	FlowerTypeID uint   `json:"flower_type_id"`
}
