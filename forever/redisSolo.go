package forever

import (
	"cloud/model"
	"cloud/utils"
	"fmt"
	"strings"
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

func GetKindsFromRedis() map[string]string {
	//get := client.Get("kinds")
	//s, e := get.Result()
	//return s, e
	all := client.HGetAll("kinds")
	m := all.Val()
	return m
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
	//client.E
	all := client.HGetAll("a")
	return all.Val()
}

func GetAllKindFromRedis() map[string]string {
	all := client.HGetAll("kinds")
	m := all.Val()
	return m
}

func CheckDocAndSave(path string) (string, error) {
	if strings.HasSuffix(path, ".doc") {
		return "", fmt.Errorf("doc format")
	}

	exists := client.HExists("content", path)
	if exists.Val() {
		get := client.HGet("content", path)
		content := get.Val()
		return content, nil
	} else {
		s, e := utils.ReadDoc(path)
		if e != nil {
			return s, e
		}
		client.HSet("content", path, s)
		return s, nil
	}
}
