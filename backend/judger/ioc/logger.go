package ioc

import (
	"5DOJ/judger/global"
	"fmt"

	"github.com/spf13/viper"
	"github.com/to404hanga/pkg404/logger"
	"go.uber.org/zap"
)

func InitLogger() {
	type Config struct {
		Level string `yaml:"level"`
	}
	var config Config
	err := viper.UnmarshalKey("Logger", &config)
	if err != nil {
		panic(fmt.Errorf("读取 logger 配置失败: %s", err))
	}

	var cfg zap.Config
	switch config.Level {
	case "debug", "dev", "test":
		cfg = zap.NewDevelopmentConfig()
	case "product":
		cfg = zap.NewProductionConfig()
	default:
		panic(fmt.Errorf("Logger.level 配置错误: %s", config.Level))
	}

	var l *zap.Logger
	l, err = cfg.Build()
	if err != nil {
		panic(fmt.Errorf("初始化 logger 失败: %s", err))
	}
	global.L = logger.NewZapLogger(l)
}
