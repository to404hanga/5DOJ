package ioc

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB() *mongo.Database {
	type Config struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"dbName"`
	}
	var config Config
	err := viper.UnmarshalKey("MongoDB", &config)
	if err != nil {
		panic(fmt.Errorf("读取 MongoDB 配置失败: %s", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 连接数据库
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", config.User, config.Password, config.Host, config.Port)))
	if err != nil {
		panic(fmt.Errorf("连接 MongoDB 失败: %s", err))
	}

	// 测试连接
	if err = client.Ping(ctx, nil); err != nil {
		panic(fmt.Errorf("Ping MongoDB 失败: %s", err))
	}

	return client.Database(config.DBName)
}
