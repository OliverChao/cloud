package forever

import (
	"cloud/model"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestBaseConRegister(t *testing.T) {
}

func TestMysqlRegister(t *testing.T) {

	MysqlRegister()
	kind := &model.Kind{
		Model: model.Model{
			ID: 3,
		},
	}
	//db.Delete(&kind)
	//db.Model(&kind).Update("deleted_at",nil)
	db.Unscoped().Model(&kind).Update("deleted_at", nil)
	MysqlUnRegister()
}
func TestMysqlFunction_1(t *testing.T) {

	MysqlRegister()
	kind := &model.Kind{
		Model: model.Model{
			ID: 3,
		},
	}
	for i := 1; i <= 10; i++ {
		go func() {
			db.Model(&kind).First(&kind)
		}()
	}
	db.Model(&kind).First(&kind)
	fmt.Println(kind.Count, kind.Name)
	MysqlUnRegister()
}
func TestMysqlFunction(t *testing.T) {

	MysqlRegister()
	//kinds := make([]*model.Kind,0)
	var kinds []*model.Kind
	db.Find(&kinds)
	fmt.Println(kinds)
	for _, v := range kinds {
		fmt.Println(v)
	}
	MysqlUnRegister()
}
func TestFunction(t *testing.T) {

	MysqlRegister()

	user := model.Admin{}
	db.Where("name = ?", "oliver").First(&user)
	fmt.Println(user.ID)

	a := make([]model.Article, 0)
	db.Offset(1*3).Limit(3).Where("user_id = ?", user.ID).Find(&a)
	//fmt.Println(a[0].ID)
	for i := range a {
		fmt.Println(a[i].ID)
	}
	MysqlUnRegister()
}

// data for mysql&redis init
func TestRedisRegister(t *testing.T) {
	MysqlDropAll()
	MysqlRegister()

	MysqlInitData()

	RedisRegister()
	RedisInitData()
	//a1 := client.HGetAll("article:1:4")
	//a2 := client.HGetAll("article:1:5")
	//a3 := client.HGetAll("article:1:6")
	//fmt.Printf(a1.String())
	//fmt.Printf(a2.String())
	//fmt.Printf(a3.String())
	RedisUnRegister()
	MysqlUnRegister()
}

func TestMysqlDropAll(t *testing.T) {
	MysqlDropAll()
	MysqlUnRegister()
}

func TestRedisGetH(t *testing.T) {
	RedisRegister()
	//b:=IsExitsKind("topic_one")
	//logrus.Info(b)
	all := client.HGetAll("船泊")
	err := all.Err()
	if err != nil {
		logrus.Error(err)
	}
	m := all.Val()
	fmt.Println(m)
	fmt.Println(len(m))
	RedisUnRegister()
}

func TestRedisHGetall(t *testing.T) {
	RedisRegister()
	v4 := uuid.NewV4()
	article := &model.Article{
		UUID:     v4.String(),
		PushedAt: time.Now(),
		Title:    "filename",
		Content:  "",
		Path:     "filepath",
		KindName: "kindname",
		HashData: "hashdata",
	}
	data := article.GenRedisData()
	bytes, _ := json.Marshal(data)

	fmt.Println(string(bytes))
	//fmt.Println(data["kind"].(string),data["id"].(string))
	//client.HSet(data["kind"].(string),data["title"].(string),string(bytes))
	m := map[string]string{}
	get := client.HGet(data["kind"].(string), data["title"].(string))
	err := json.Unmarshal([]byte(get.Val()), &m)
	if err != nil {
		logrus.Error(err)
		return
	}
	fmt.Println(m)

	RedisUnRegister()
}

func TestMysqlInitData(t *testing.T) {
	MysqlRegister()
	MysqlUnRegister()
	//MysqlInitData()
}
