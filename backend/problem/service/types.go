package service

import (
	"5DOJ/problem/domain"
	"context"
)

type IProblemService interface {
	Get(ctx context.Context, pid uint64) (problemView domain.ProblemView, err error)
	GetTestCaseList(ctx context.Context, pid uint64) (testCaseList []domain.TestCaseView, err error)
	GetList(ctx context.Context, size int, cursorIn uint64, title string) (cursorOut uint64, list []domain.ProblemView, err error)
	Create(ctx context.Context, title string, level int, createdBy uint64, timeLimit, memoryLimit int, markdown string) (pid uint64, err error)
	Update(ctx context.Context, pid uint64, title string, level int, updatedBy uint64, timeLimit, memoryLimit int, markdown string) (err error)
	Enable(ctx context.Context, pid, updatedBy uint64) (err error)
	Disable(ctx context.Context, pid, updatedBy uint64) (err error)
	AppendTestCase(ctx context.Context, pid uint64, input, output string, score int, createdBy uint64) (tid string, err error)
	UpdateTestCase(ctx context.Context, pid uint64, tid, input, output string, score int, updatedBy uint64) (err error)
	EnableTestCase(ctx context.Context, pid uint64, tid string, updatedBy uint64) (err error)
	DisableTestCase(ctx context.Context, pid uint64, tid string, updatedBy uint64) (err error)
}
