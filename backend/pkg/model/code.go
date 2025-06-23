package model

import (
	"5DOJ/submitter/global"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Code struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	FilenameWithoutExt string             `bson:"filenameWithoutExt"`
	Content            string             `bson:"content"`
}

func InsertCode(ctx context.Context, code Code) (err error) {
	_, err = global.MongoDB.Collection("code").InsertOne(ctx, &code)
	return err
}
