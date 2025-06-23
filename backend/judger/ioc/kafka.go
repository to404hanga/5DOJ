package ioc

import (
	"5DOJ/judger/consumer"
	"5DOJ/judger/global"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"github.com/to404hanga/pkg404/saramax"
)

func InitKafka() {
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

	global.Kafka, err = sarama.NewClient(config.Addrs, saramaCfg)
	if err != nil {
		panic(fmt.Errorf("连接 Kafka 失败: %s", err))
	}
}

func NewConsumers(judger *consumer.JudgerSubmitConsumer) []saramax.Consumer {
	return []saramax.Consumer{
		judger,
	}
}
