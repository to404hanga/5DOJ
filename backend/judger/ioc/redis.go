package ioc

import (
	"5DOJ/judger/global"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() {
	type Config struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
	var config Config
	if err := viper.UnmarshalKey("Redis", &config); err != nil {
		panic(err)
	}

	global.Rds = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
}
