package forever

import (
	"cloud/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ego/riot/types"
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

func AddFileToKind(kindname, filepath, title, hashdata string, count int) {
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
			Title:    title,
			Content:  "",
			Path:     filepath,
			KindName: kindname,
			HashData: hashdata,
			FullName: fullName,
		}
		//db.Create(&article)
		// 先入mysql 会更新 article的值, 补全后,再存入redis
		_ = addArticleMutex(article)
		addArticleInfoToRedis(article)
		// update search engine
		AddDoc(kindname, title)
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

func GetOneArticleFromRedisByTitle(kindName, title string) string {
	get := client.HGet(kindName, title)
	m := map[string]string{}
	_ = json.Unmarshal([]byte(get.Val()), &m)
	logrus.Info(m["path"])
	return m["path"]
}

func GetAllArticleFromRedis(kindName string) (map[string]string, error) {
	all := client.HGetAll(kindName)
	m := all.Val()
	if len(m) == 0 {
		return m, errors.New("kind name no exists")
	}
	for k, _ := range m {
		m[k] = kindName
	}
	return m, nil
}

func GetSearchResultFromRiot(docs []types.ScoredDoc) map[string]string {
	m := map[string]string{}
	for _, v := range docs {
		m[v.Content] = v.Attri.(string)
	}
	return m
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

// 删除resource/ 下的文件, 更改数据库
func DeleteArticleFunc(name, kind string) error {
	get := client.HGet(kind, name)
	s, e := get.Result()
	if e != nil {
		//logrus.Error("redis no this data")
		return e
	}
	// by is the number after change redis
	by := client.HIncrBy("kinds", kind, -1)
	num := by.Val()
	//logrus.Info(by.Val())
	m := map[string]string{}
	_ = json.Unmarshal([]byte(s), &m)
	logrus.Info(m)
	article := &model.Article{
		Title:    m["title"],
		KindName: m["kind"],
		HashData: m["hash"],
	}
	k := &model.Kind{
		Name: m["kind"],
	}
	db.Where(&article).Find(&article).Delete(&article)
	db.Model(&model.Kind{}).Where(&k).Select("count").Update("count", int(num))
	logrus.Info(article)
	err := os.Remove(m["path"])
	if err != nil {
		logrus.Error("[DeleteArticleFunc]", err)
	}
	//db.Model(&article).Delete(&article)
	client.HDel(kind, name)
	// 清除riot引擎索引
	RiotDeleteDocByNameKind(kind, name)
	return nil
}

func DeleteKindFunc(kind string) error {

	get := client.HDel("kinds", kind)
	i, e := get.Result()
	if e != nil {
		logrus.Error(e)
		return e
	}
	logrus.Info(i)
	if i == 0 {
		return fmt.Errorf("no kind %s", kind)
	}
	//清除 redis 数据
	client.Expire(kind, 0)

	db.Where("kind_name like ?", kind).Delete(model.Article{})
	db.Where("name like ?", kind).Delete(model.Kind{})

	path := "resource/" + kind
	err := os.RemoveAll(path)
	if err != nil {
		logrus.Error(err)
	}
	return nil
}
