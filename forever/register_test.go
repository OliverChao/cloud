package forever

import (
	"cloud/model"
	"encoding/json"
	"fmt"
	"testing"
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

	//InitResourceDirs()

	RedisUnRegister()
	MysqlUnRegister()
}

func TestMysqlDropAll(t *testing.T) {
	//MysqlDropAll()
	MysqlRegister()
	a := &model.Admin{
		Name: "admin",
	}
	db.Model(&a).Where("name = ?", a.Name).Update("password","qwe")
	MysqlUnRegister()
}

func TestRedisGetH(t *testing.T) {
	RedisRegister()
	//b:=IsExitsKind("topic_one")
	//logrus.Info(b)
	kind := "国际公约"
	name := "t"
	get := client.HGet(kind, name)
	//s, _ := get.Result()
	s := get.Val()
	//fmt.Println(get.Result())
	m := map[string]string{}
	json.Unmarshal([]byte(s), &m)
	fmt.Println(m)
	fmt.Println(m["hash"])
	fmt.Println(m["title"])

	RedisUnRegister()
}

func TestRedisHGetall(t *testing.T) {
	RedisRegister()
	InitResourceDirs()
	//v4 := uuid.NewV4()
	//article := &model.Article{
	//	UUID:     v4.String(),
	//	PushedAt: time.Now(),
	//	Title:    "filename",
	//	Content:  "",
	//	Path:     "filepath",
	//	KindName: "kindname",
	//	HashData: "hashdata",
	//}
	//data := article.GenRedisData()
	//bytes, _ := json.Marshal(data)
	//
	//fmt.Println(string(bytes))
	////fmt.Println(data["kind"].(string),data["id"].(string))
	////client.HSet(data["kind"].(string),data["title"].(string),string(bytes))
	//m := map[string]string{}
	//get := client.HGet(data["kind"].(string), data["title"].(string))
	//err := json.Unmarshal([]byte(get.Val()), &m)
	//if err != nil {
	//	logrus.Error(err)
	//	return
	//}
	//fmt.Println(m)

	RedisUnRegister()
}

func TestMysqlInitData(t *testing.T) {
	MysqlRegister()
	var articleList []model.Article
	db.Order("created_at desc,id").Limit(10).Find(&articleList)
	//fmt.Println(articleList)
	for _, v := range articleList {
		fmt.Println(v.FullName)
	}
	MysqlUnRegister()
	//MysqlInitData()
}
