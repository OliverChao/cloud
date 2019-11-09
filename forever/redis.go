package forever

import (
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"cloud/config/redisCon"
)

var client *redis.Client

func ConnectRedis(con *redisCon.RedisConfig) {
	//client := redis.NewClient{con.RedisOption}
	client = redis.NewClient(con.RedisOption)
	logrus.Info("[redis]connect successfully...")
}

func DisConnectRedis() {
	if e := client.Close(); e != nil {
		logrus.Error("[redis]close failed...")
	} else {
		logrus.Info("[redis]close successfully")
	}
}

func GetGlobalRedisClient() *redis.Client {
	return client
}
