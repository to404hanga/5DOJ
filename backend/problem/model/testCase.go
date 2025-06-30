package model

type TestCase struct {
	Tid       string `bson:"tid"`
	Pid       uint64 `bson:"pid"`
	Input     string `bson:"input"`
	Output    string `bson:"output"`
	Score     int    `bson:"score"`
	CreatedBy string `bson:"createdBy"`
	UpdatedBy string `bson:"updatedBy"`
	Enabled   bool   `bson:"enabled"`
}

func (TestCase) TableName() string {
	return "test_case"
}
