package baseCon

import (
	"github.com/sirupsen/logrus"
)

const (
	Version = "0.0.1"
)

type BaseConfig struct { //	redisOp := &redis.Options{
	//		Addr:     "127.0.0.1:6379",
	//		Password: "toor",
	//		DB:       0,
	//		PoolSize: 10,
	//	}
	Host     string
	Port     int
	LogLevel logrus.Level
}

func LoadBaseConfig() *BaseConfig {
	IFcon := &BaseConfig{
		Host: "127.0.0.1",
		Port: 8080,
		//MysqlUri: mysqlSign,
		//RedisUri: redisOp,
		LogLevel: logrus.InfoLevel,
	}
	return IFcon
}

//func init() {
//	mysqlSign := "root:toor@(127.0.0.1:3306)/ifei?charset=utf8mb4&parseTime=True&loc=Local"

//
//	IFcon := &BaseConfig{
//		Host:     "127.0.0.1",
//		Port:     8080,
//		//MysqlUri: mysqlSign,
//		//RedisUri: redisOp,
//		LogLevel: logrus.InfoLevel,
//	}
//
//}
