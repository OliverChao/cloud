package forever

import (
	"cloud/config/baseCon"
	"cloud/config/mysqlCon"
	"cloud/config/redisCon"
	"cloud/model"
	"github.com/sirupsen/logrus"
	"os"
)

func BaseConRegister() {
	IFconBase := baseCon.LoadBaseConfig()
	logrus.SetLevel(IFconBase.LogLevel)
}

func MysqlRegister() {
	IFMysqlCon := mysqlCon.LoadMysqlConfig()
	Connect(IFMysqlCon)
}

func MysqlDropAll() {
	IFMysqlCon := mysqlCon.LoadMysqlConfig()
	Connect(IFMysqlCon)
	DropAll(IFMysqlCon)
}

func RedisRegister() {
	redisConfig := redisCon.LoadRedisConfig()
	ConnectRedis(redisConfig)
	//RedisInitData()
}

func MysqlInitData() {
	kinds := []*model.Kind{
		{
			Name:  "国际公约",
			Count: 0,
		},
		{
			Name:  "法律",
			Count: 0,
		},
		{
			Name:  "行政法规",
			Count: 0,
		},
		{
			Name:  "部门规章",
			Count: 0,
		},
	}
	for _, v := range kinds {
		db.Create(&v)
	}

	//pwd : Sealaw@520
	admin := &model.Admin{
		Name:              "admin",
		Password:          "a0404ae3334ec51621e5161a13e82caa320d046bef7fd1f8b34e355d1968beef",
		TotalArticleCount: 0,
	}
	db.Create(&admin)

}

func RedisInitData() {
	client.FlushAll()

	var kinds []*model.Kind
	db.Find(&kinds)
	for _, v := range kinds {
		client.HSet("kinds", v.Name, v.Count)
	}
	//db.Find()
}

func InitResourceDirs() {
	all := client.HGetAll("kinds")
	m := all.Val()
	for k, _ := range m {
		logrus.Info(k)
		err := os.Mkdir("resource/"+k, os.ModePerm)
		if err != nil {
			logrus.Error(err)
		}

	}
}

//func DropResourceDirs() {
//	all :=
//}
