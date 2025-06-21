package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/to404hanga/5DOJ/pkg/constant"
	"github.com/to404hanga/5DOJ/program/internal/global"
	"github.com/to404hanga/5DOJ/program/internal/model"
)

type ProgramService struct {
}

var _ IProgramService = (*ProgramService)(nil)

func NewProgramService() *ProgramService {
	return &ProgramService{}
}

func (p *ProgramService) CreateProgram(ctx context.Context, title, content string, createdBy uint64, level constant.ProgramLevelType, timeLimitMS, memoryLimitMB uint64, tags []string) (programId uint64, err error) {
	now := time.Now().UTC()

	program := &model.ProgramBase{
		Title:       title,
		Level:       int8(level),
		Tags:        strings.Join(tags, ","),
		TimeLimit:   timeLimitMS * 1000000,
		MemoryLimit: memoryLimitMB * 1024 * 1024,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = global.SqlDB.WithContext(ctx).Create(&program).Error
	return program.Id, err
}

func (p *ProgramService) AppendProgramTestCase(ctx context.Context, programId, appendedBy uint64, input, expected string) (testCaseNum int, testCaseId string, err error) {
	now := time.Now().UTC()

	var cnt int
	if err = global.SqlDB.WithContext(ctx).
		Where("id = ?", programId).
		Select("test_case_num").
		Scan(&cnt).Error; err != nil {
		return
	}
	testCaseNum = cnt + 1
	testCaseId = fmt.Sprintf("P%d:T%d:%s", programId, testCaseNum, now.Format(time.DateOnly))

	if err = model.AppendOneTestCase(ctx, programId, testCaseId, input, expected); err != nil {
		return 0, "", err
	}

	if err = global.SqlDB.WithContext(ctx).
		Model(&model.ProgramBase{}).
		Where("id = ?", programId).
		Update("test_case_num", testCaseNum).Error; err != nil {
		// 异步重试删除，确保数据一致性
		go func() {
			for {
				if err := model.DeleteOneTestCaseByProgramIdAndTestCaseId(context.Background(), programId, testCaseId); err == nil {
					return
				}
			}
		}()
		return 0, "", err
	}

	return testCaseNum, testCaseId, nil
}
