package global

import (
	"github.com/IBM/sarama"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	MySQL   *gorm.DB
	MongoDB *mongo.Database
	Kafka   sarama.Client
)
