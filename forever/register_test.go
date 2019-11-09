package forever

import (
	"fmt"
	"iFei/model"
	"testing"
)

func TestBaseConRegister(t *testing.T) {
}

func TestMysqlRegister(t *testing.T) {

	MysqlRegister()

	MysqlUnRegister()
}
func TestMysqlFunction(t *testing.T) {

	MysqlRegister()
	//QueryDemo()
	//MysqlDropAll()
	//CreateDemo()
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
	RedisRegister()
	a1 := client.HGetAll("article:1:4")
	a2 := client.HGetAll("article:1:5")
	a3 := client.HGetAll("article:1:6")
	fmt.Printf(a1.String())
	fmt.Printf(a2.String())
	fmt.Printf(a3.String())
	RedisUnRegister()
}

func TestMysqlDropAll(t *testing.T) {
	MysqlDropAll()
	MysqlUnRegister()
}
