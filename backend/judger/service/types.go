package service

import (
	"5DOJ/pkg/constant/contestMode"
	"5DOJ/pkg/constant/evaluationStatus"
	"5DOJ/pkg/constant/language"
	"context"
)

type IJudgerService interface {
	Preheater(ctx context.Context, contestId uint64) (err error)
	Judge(ctx context.Context, recordId, problemId uint64, lang language.LanguageType, filenameWithoutExt, userCode string, mode contestMode.ContestModeType) (evalutionStatus evaluationStatus.EvaluationStatusType, timeUsageMS, memoryUsageKB uint64, err error)
}
