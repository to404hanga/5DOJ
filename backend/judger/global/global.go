package global

import (
	"5DOJ/judger/domain"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"github.com/to404hanga/pkg404/logger"
	"gorm.io/gorm"
)

var (
	MySQL *gorm.DB
	CP    map[uint64]domain.Problem
	Rds   redis.Cmdable
	L     logger.Logger
	Kafka sarama.Client
)
