package model

import (
	"context"
	"time"

	"github.com/to404hanga/5DOJ/program/internal/global"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestCase struct {
	Id       string `bson:"id"`       // 测试用例 ID
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

func GetProgramContentByProgramId(ctx context.Context, programId uint64) (ProgramContent, error) {
	var program ProgramContent
	err := global.MongoDB.Collection("program").FindOne(ctx, bson.M{"program_id": programId}).Decode(&program)
	return program, err
}

func UpdateOneTestCaseByProgramIdAndTestCaseId(ctx context.Context, programId uint64, testCaseId, input, expected string) error {
	update := bson.M{
		"$set": bson.M{
			"test_cases.$.input":    input,
			"test_cases.$.expected": expected,
		},
	}
	filter := bson.M{
		"program_id":    programId,
		"test_cases.id": testCaseId,
	}

	_, err := global.MongoDB.Collection("program").UpdateOne(ctx, filter, update)
	return err
}

func AppendOneTestCase(ctx context.Context, programId uint64, testCaseId, input, expected string) error {
	update := bson.M{
		"$push": bson.M{
			"test_cases": bson.M{
				"id":       testCaseId,
				"input":    input,
				"expected": expected,
			},
		},
	}
	filter := bson.M{
		"program_id": programId,
	}

	_, err := global.MongoDB.Collection("program").UpdateOne(ctx, filter, update)
	return err
}

func DeleteOneTestCaseByProgramIdAndTestCaseId(ctx context.Context, programId uint64, testCaseId string) error {
	update := bson.M{
		"$pull": bson.M{
			"test_cases": bson.M{
				"id": testCaseId,
			},
		},
	}
	filter := bson.M{
		"program_id": programId,
	}

	_, err := global.MongoDB.Collection("program").UpdateOne(ctx, filter, update)
	return err
}
