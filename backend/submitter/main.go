package main

import (
	"5DOJ/submitter/global"
	"5DOJ/submitter/ioc"
	"5DOJ/submitter/producer"
	"5DOJ/submitter/rpc"
	"5DOJ/submitter/service"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	cfile := pflag.String("config", "config/config.yaml", "配置文件路径")
	pflag.Parse()

	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %s", err))
	}

	global.MySQL = ioc.InitMySQL()
	global.MongoDB = ioc.InitMongoDB()
	global.L = ioc.InitLogger()
	global.Kafka = ioc.InitKafka()
	global.Etcd = ioc.InitEtcd()
	global.GrpcServer = ioc.InitGrpcServer(rpc.NewSubmitterServiceServer(service.NewSubmitterService(producer.NewSubmitterProducer())))
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("Prometheus.port")), nil)
	}()

	if err := global.GrpcServer.Serve(); err != nil {
		panic(fmt.Errorf("启动Grpc服务失败: %s", err))
	}
}
