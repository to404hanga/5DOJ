package user

import (
	"5DOJ/user/global"
	"5DOJ/user/ioc"
	"5DOJ/user/rpc"
	"5DOJ/user/service"
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
	global.L = ioc.InitLogger()
	global.Etcd = ioc.InitEtcd()
	global.GrpcServer = ioc.InitGrpcServer(rpc.NewUserServiceServer(service.NewUserService(viper.GetString("app.defaultPassword"))))
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("Prometheus.port")), nil)
	}()

	if err := global.GrpcServer.Serve(); err != nil {
		panic(fmt.Errorf("启动 Grpc 服务失败: %s", err))
	}
}
