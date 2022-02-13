package pika

import (
	viper2 "forever.love/initialize/viper"
	"github.com/go-redis/redis"
)

var (
	Client *redis.Client
)

func Setup() {
	Client = redis.NewClient(&redis.Options{
		Addr:     viper2.RedisConf.Host,
		Password: viper2.RedisConf.Password,
		DB:       viper2.RedisConf.DBNum,
	})
}
