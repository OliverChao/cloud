package forever

import (
	"cloud/model"
	"fmt"
	"github.com/sirupsen/logrus"
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

func TestRedisRegister(t *testing.T) {
	MysqlRegister()
	RedisRegister()
	//a1 := client.HGetAll("article:1:4")
	//a2 := client.HGetAll("article:1:5")
	//a3 := client.HGetAll("article:1:6")
	//fmt.Printf(a1.String())
	//fmt.Printf(a2.String())
	//fmt.Printf(a3.String())
	RedisInitData()
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
	all := client.HGetAll("kinds")
	val := all.Val()
	for k, v := range val {
		fmt.Println(k, " => ", v)
	}
	RedisUnRegister()
}

func TestRedisHGetall(t *testing.T) {
	RedisRegister()
	//b:=IsExitsKind("topic_one")
	//logrus.Info(b)
	m := test()
	logrus.Info(m)
	fmt.Println(m["1"])
	RedisUnRegister()
}

func TestMysqlInitData(t *testing.T) {
	MysqlRegister()
	MysqlUnRegister()
	//MysqlInitData()
}
