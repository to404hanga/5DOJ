package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Code struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	FilenameWithoutExt string             `bson:"filenameWithoutExt"`
	Content            string             `bson:"content"`
}
