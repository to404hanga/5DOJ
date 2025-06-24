package main

import (
	"5DOJ/judger/consumer"
	"5DOJ/judger/global"
	"5DOJ/judger/ioc"
	"5DOJ/judger/service"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	global.Kafka = ioc.InitKafka()
	global.Redis = ioc.InitRedis()
	global.L = ioc.InitLogger()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("Prometheus.port")), nil)
	}()

	goJudgeSvc := service.NewGoJudgeJudgerService(viper.GetString("GoJudge.baseUrl"))
	judgerConsumer := consumer.NewJudgerSubmitConsumer(goJudgeSvc)

	consumers := ioc.NewConsumers(judgerConsumer)

	global.L.Info("开始启动消费者")
	for _, c := range consumers {
		if err := c.Start(); err != nil {
			panic(fmt.Errorf("启动消费者失败: %s", err))
		}
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
}
