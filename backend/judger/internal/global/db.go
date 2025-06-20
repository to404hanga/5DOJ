package global

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	SqlDB   *gorm.DB
	MongoDB *mongo.Database
)
