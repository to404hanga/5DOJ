package ioc

import (
	"fmt"

	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InitEtcd() *clientv3.Client {
	var config clientv3.Config
	err := viper.UnmarshalKey("Etcd", &config)
	if err != nil {
		panic(fmt.Errorf("读取 Etcd 配置失败: %s", err))
	}

	client, err := clientv3.New(config)
	if err != nil {
		panic(fmt.Errorf("连接 Etcd 失败: %s", err))
	}
	return client
}
