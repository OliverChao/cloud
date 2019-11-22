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

func GetKindsFromRedis() (string, error) {
	get := client.Get("kinds")
	s, e := get.Result()
	return s, e
}

func IsExitsKind(name string) bool {
	//get := client.HGet("kinds", name)
	exists := client.HExists("kinds", name)
	//logrus.Info(exists.Val())
	//logrus.Info(exists.String())
	return exists.Val()
}

func HSet(name, count string) {
	if err := client.HSet("kinds", name, count).Err(); err != nil {
	}
}

func AddKindToRedis(name string) {
	client.HSet("kinds", name, 0)
}

func test() map[string]string {
	all := client.HGetAll("a")
	return all.Val()
}

func GetAllKindFromRedis() map[string]string {
	all := client.HGetAll("kinds")
	m := all.Val()
	return m
}
