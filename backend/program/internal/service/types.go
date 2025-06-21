package service

import (
	"context"

	"github.com/to404hanga/5DOJ/pkg/constant"
	"github.com/to404hanga/5DOJ/program/internal/domain"
)

type IProgramService interface {
	CreateProgram(ctx context.Context, title, content string, createdBy uint64, level constant.ProgramLevelType, timeLimitMS, memoryLimitMB uint64, tags []string) (programId uint64, err error)
	AppendProgramTestCase(ctx context.Context, programId, appendedBy uint64, input, expected string) (testCaseNum int, testCaseId string, err error)
	UpdateTestCaseByProgramIdAndTestCaseId(ctx context.Context, programId uint64, testCaseId, intput, expected string, updatedBy uint64) (err error)
	DeleteTestCaseByProgramIdAndTestCaseId(ctx context.Context, programId uint64, testCaseId string, deletedBy uint64) (err error)
	GetProgramByProgramId(ctx context.Context, programId uint64) (program domain.ProgramView, err error)
	GetTestCasesByProgramId(ctx context.Context, programId uint64) (testCases []domain.TestCaseView, err error)
	EnableProgram(ctx context.Context, programId uint64, enabledBy uint64) (err error)
	DisableProgram(ctx context.Context, programId uint64, disabledBy uint64) (err error)
}
