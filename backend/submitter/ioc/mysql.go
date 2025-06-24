package ioc

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

func InitMySQL() *gorm.DB {
	type Config struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"dbName"`
	}
	var config Config
	err := viper.UnmarshalKey("MySQL", &config)
	if err != nil {
		panic(fmt.Errorf("读取 MySQL 配置失败: %s", err))
	}

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=UTC", config.User, config.Password, config.Host, config.Port, config.DBName)), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("连接 MySQL 失败: %s", err))
	}

	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          config.DBName,
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		},
	}))
	if err != nil {
		panic(fmt.Errorf("注册 MySQL 监控失败: %s", err))
	}

	return db
}
