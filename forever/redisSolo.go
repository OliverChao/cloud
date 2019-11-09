package forever

import (
	"cloud/model"
)

func GetArticleInfo() string {
	cmd := client.HGet("user:1:info", "age")
	s, e := cmd.Result()
	//return fmt.Sprintf("%v",cmd.Result())
	if e != nil {
		return ""
	}
	return s
}

func StoreArticleInfo(article *model.Article) {
	//m := article.RetArticle()
	//key := fmt.Sprintf("article:%v:%v", article.UserID, article.ID)
	//client.HMSet(key, m)
	//client.Expire(key, 11400*time.Second)
}
