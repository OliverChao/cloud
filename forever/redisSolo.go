package forever

import (
	"cloud/model"
	"cloud/utils"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
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

// func GetArticleNumberByRedis()map[string]string{
// 	get := client.HGetAll("kinds")
// 	return get.Val()
// }

func CheckDocAndSave(path string) (string, error) {
	if strings.HasSuffix(path, ".doc") {
		return "", fmt.Errorf("doc format")
	}

	if !client.HExists("size", path).Val() {
		return "", fmt.Errorf("can't get file size")
	} else {
		sizeCmd := client.HGet("size", path)
		val := sizeCmd.Val()
		//strconv.Atoi()
		size, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return "", err
		}
		if strings.HasSuffix(path, ".docx") {
			// 200K
			if size > 204800 {
				logrus.Info(path, "is above 200K")
				return "", fmt.Errorf("big docx file")
			}
		}
		if strings.HasSuffix(path, ".pdf") {
			// 1M
			if size > 204800 {
				logrus.Info(path, "is above 1M")
				return "", fmt.Errorf("big pdf file")
			}
		}
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

func SaveSizeWhenUpload(path string, size int64) {
	client.HSet("size", path, size)
}

//> 1. 哈希一样, 名字一样 => 直接 卡掉, 不执行操作.                2
//> 2. 哈希一样, 名字不一样 => 正常上传                           0
//> 3. 哈希不一样, 名字一样 =>  先删除原来的,文件, 后正常上传,       1
//> 4. 哈希与名字均不一样 => 直接上传即可                          0
//>
func CheckDoubleUploadFiles(kind, title, hashData string) int {
	get := client.HGet(kind, title)
	s, err := get.Result()
	if err != nil {
		return 0
	}
	m := map[string]string{}
	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		return 0
	}

	h := m["hash"]
	t := m["title"]
	if h != hashData && t != title {
		return 0
	} else if h != hashData && t == title {
		return 1
	} else if h == hashData && t == title {
		return 2
	} else {
		return -1
	}

}
