package model

type Kind struct {
	Model
	Name  string `json:"type_name" gorm:"unique"`
	Count int    `json:"count" gorm:"default:0"`
}
