package service

import (
	"context"

	"github.com/to404hanga/5DOJ/judger/internal/domain"
)

type IJudgerService interface {
	Compile(ctx context.Context, compiler domain.Compiler, filename string, content string) (fileId string, err error)
	Run(ctx context.Context, filenameWithoutExtension, fileId, input string, timeLimit, memoryLimit int) (output string, err error)
	Judge(ctx context.Context, fileId string, programId string) (err error)
	Delete(ctx context.Context, fileId string) (err error)
}
