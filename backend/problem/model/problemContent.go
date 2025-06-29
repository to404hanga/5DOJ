package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProblemContent struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Pid      uint64             `bson:"pid"`
	Markdown string             `bson:"markdown"`
}

func (ProblemContent) TableName() string {
	return "problem_content"
}
