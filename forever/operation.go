package forever

import (
	"cloud/model"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mu *sync.Mutex = &sync.Mutex{}

func AddNewKind(name string) bool {
	isExitsKind := IsExitsKind(name)
	if !isExitsKind {
		err := os.Mkdir("resource/"+name, os.ModePerm)
		if err != nil {
			logrus.Error(err)
			return false
		}
	}
	AddKindToMysql(name)
	AddKindToRedis(name)
	return true
}

func AddFileToKind(kindname, filepath, filename, hashdata string, count int) {
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	splits := strings.Split(filepath, "/")
	fullName := splits[len(splits)-1]
	wg.Add(count)
	go func() {
		defer wg.Done()
		v4 := uuid.NewV4()
		article := &model.Article{
			UUID:     v4.String(),
			PushedAt: time.Now(),
			Title:    filename,
			Content:  "",
			Path:     filepath,
			KindName: kindname,
			HashData: hashdata,
			FullName: fullName,
		}
		//db.Create(&article)
		addArticleInfoToRedis(article)
		_ = addArticleMutex(article)
	}()

	kind := &model.Kind{
		Name: kindname,
	}
	//mu.Lock()
	//defer mu.Unlock()
	get := client.HGet("kinds", kindname)
	val := get.Val()
	i, _ := strconv.Atoi(val)
	changeNum := i + count
	client.HSet("kinds", kindname, changeNum)
	db.Model(&kind).Where("name like ?", kindname).Update("count", changeNum)
}

func addArticleInfoToRedis(a *model.Article) {
	data := a.GenRedisData()
	bytes_, _ := json.Marshal(data)
	// if error ???
	client.HSet(data["kind"].(string), data["title"].(string), string(bytes_))
}

func GetArticleFromRedisByTitle(kindName, title string) map[string]string {
	m := map[string]string{}
	get := client.HGet(kindName, title)
	_ = json.Unmarshal([]byte(get.Val()), &m)
	return m
}

func GetAllArticleFromRedis(kindName string) (map[string]string, error) {
	all := client.HGetAll(kindName)
	m := all.Val()
	if len(m) == 0 {
		return m, errors.New("kind name no exists")
	}
	return m, nil
}

func addArticleMutex(a *model.Article) (err error) {
	mu.Lock()
	defer mu.Unlock()
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	if err = tx.Create(&a).Error; nil != err {
		fmt.Println(err)
		return
	}
	return nil
}
