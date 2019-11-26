package model

import (
	"strconv"
	"time"
)

type Article struct {
	Model
	UUID     string    `json:"uuid" gorm:"unique"`
	PushedAt time.Time `json:"pushed_at"`
	Title    string    `gorm:"size:255" json:"title"`
	FullName string    `gorm:"size:255" json:"filename"`
	Content  string    `gorm:"type:text" json:"content"`
	Path     string    `gorm:"size:255" json:"path"`
	KindName string    `json:"kind_name" sql:"index"`
	HashData string    `json:"hash"`
}

func (t *Article) GenRedisData() map[string]interface{} {
	m := map[string]interface{}{}
	m["id"] = strconv.Itoa(int(t.ID))
	m["title"] = t.Title
	m["path"] = t.Path
	m["kind"] = t.KindName
	m["hash"] = t.HashData
	m["filename"] = t.FullName
	return m
}
