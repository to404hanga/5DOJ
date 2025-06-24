package ioc

import (
	"5DOJ/submitter/global"
	"5DOJ/submitter/rpc"
	"fmt"

	"github.com/spf13/viper"
	"github.com/to404hanga/pkg404/grpcx"
	"google.golang.org/grpc"
)

// 需要先初始化 global.Etcd 和 global.L
func InitGrpcServer(submitter *rpc.SubmitterServiceServer) *grpcx.Server {
	type Config struct {
		Port     int    `yaml:"port"`
		EtcdAddr string `yaml:"etcdAddr"`
		EtcdTTL  int64  `yaml:"etcdTTL"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.server", &config)
	if err != nil {
		panic(fmt.Errorf("读取 grpc server 配置失败: %s", err))
	}

	srv := grpc.NewServer()
	submitter.Register(srv)

	return &grpcx.Server{
		Server:     srv,
		Port:       config.Port,
		EtcdClient: global.Etcd,
		Name:       "submitter",
		EtcdTTL:    config.EtcdTTL,
		L:          global.L,
	}
}
