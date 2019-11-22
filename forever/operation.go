package forever

import (
	"cloud/model"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
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

func AddFileToKind(kindname, filepath, filename, hashdata string) {

	v4 := uuid.NewV4()
	article := &model.Article{
		UUID:     v4.String(),
		PushedAt: time.Now(),
		Title:    filename,
		Content:  "",
		Path:     filepath,
		KindName: kindname,
		HashData: hashdata,
	}
	db.Create(&article)

	kind := &model.Kind{
		Name: kindname,
	}
	get := client.HGet("kinds", kindname)
	val := get.Val()
	i, _ := strconv.Atoi(val)
	client.HSet("kinds", kindname, i+1)
	db.Model(&kind).Where("name like ?", kindname).Update("count", i+1)
}
