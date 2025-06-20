package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestCase struct {
	Input    string `bson:"input"`    // 测试用例输入
	Expected string `bson:"expected"` // 测试用例对应答案
}

type ProgramContent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // 自动生成 ID
	ProgramId uint64             `bson:"program_id"`    // 所属题目 ID
	Content   string             `bson:"content"`       // 题目正文
	TestCases []TestCase         `bson:"test_cases"`    // 测试用例数组
	CreatedAt time.Time          `bson:"created_at"`    // 创建时间
	UpdatedAt time.Time          `bson:"updated_at"`    // 更新时间
}
