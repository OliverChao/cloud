package redisCon

import "github.com/go-redis/redis/v7"

var (
	RedisOp = &redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "toor",
		DB:       0,
		PoolSize: 10,
	}
)

//redis config like mysql config struct
type RedisConfig struct {
	RedisOption *redis.Options
}

func LoadRedisConfig() (redisConfig *RedisConfig) {
	redisConfig = &RedisConfig{RedisOption: RedisOp}

	return
}
