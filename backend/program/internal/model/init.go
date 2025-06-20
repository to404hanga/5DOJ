package model

import (
	"context"
	"time"

	"github.com/to404hanga/5DOJ/program/internal/global"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMySQL() error {
	return global.SqlDB.AutoMigrate(
		&ProgramBase{},
	)
}

func InitMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	col := global.MongoDB.Collection("program")
	if _, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"program_id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"test_cases.id": "text"},
			Options: options.Index().SetUnique(true),
		},
	}); err != nil {
		return err
	}

	return nil
}
