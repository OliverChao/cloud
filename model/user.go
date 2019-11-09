package model

type Admin struct {
	Model
	Name              string `sql:"index" gorm:"size:255;unique" json:"name"`
	Password          string `gorm:"not null"`
	TotalArticleCount int    `json:"total_article_count"`

}

func (u *Admin) RetData() map[string]interface{} {
	return nil
}
