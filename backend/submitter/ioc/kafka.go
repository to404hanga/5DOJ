package ioc

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var config Config
	err := viper.UnmarshalKey("Kafka", &config)
	if err != nil {
		panic(fmt.Errorf("读取 Kafka 配置失败: %s", err))
	}

	client, err := sarama.NewClient(config.Addrs, saramaCfg)
	if err != nil {
		panic(fmt.Errorf("连接 Kafka 失败: %s", err))
	}
	return client
}
