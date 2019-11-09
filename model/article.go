package model

import "time"

type Article struct {
	Model
	UUID     string    `json:"uuid" gorm:"unique"`
	PushedAt time.Time `json:"pushed_at"`
	Title    string    `gorm:"size:255" json:"title"`
	Content  string    `gorm:"type:text" json:"content"`
	Path     string    `sql:"index" gorm:"size:255" json:"path"`
}