package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TestCase struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Pid       uint64             `bson:"pid"`
	Input     string             `bson:"input"`
	Output    string             `bson:"output"`
	Score     int                `bson:"score"`
	CreatedBy uint64             `bson:"createdBy"`
	UpdatedBy uint64             `bson:"updatedBy"`
	Enabled   bool               `bson:"enabled"`
}
