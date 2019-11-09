package forever

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"cloud/config/redisCon"
	"cloud/model"
	"testing"
)

func TestConnectionRedis(t *testing.T) {
	redisConfig := redisCon.LoadRedisConfig()

	ConnectRedis(redisConfig)
	client.FlushAll()
	client.HSet("user:1:info", "age", 20)
	get := client.HGet("user:1:info", "age")
	s, _ := get.Result()
	if s != "20" {
		t.Fail()
	}
	DisConnectRedis()
}

func TestGetRedisMsg(t *testing.T) {
	redisConfig := redisCon.LoadRedisConfig()

	ConnectRedis(redisConfig)

	//exists := client.HExists("user:oliver:info", "articles")
	//fmt.Println(exists.Result())
	get := client.HGet("user:oliver:info", "articles")
	bytes, _ := get.Bytes()
	//
	////fmt.Println(bytes)
	var a []model.Article
	e := json.Unmarshal(bytes, &a)
	if e != nil {
		logrus.Error(e)
	} else {
		//fmt.Println(a)
	}
	//
	for i, v := range a {
		fmt.Println(i, v.Content)
	}
	DisConnectRedis()
}
