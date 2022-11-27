package models

type FlowerWithType struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	FlowerTypeID   uint   `json:"flower_type_id"`
	FlowerTypeName string `json:"flower_type_name"`
}
