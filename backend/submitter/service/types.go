package service

import (
	"5DOJ/submitter/domain"
	"context"
)

type ISubmitterService interface {
	Submit(ctx context.Context, contestId, problemId, userId uint64, lang, code, mode string) (recordId uint64, err error)
	Query(ctx context.Context, recordId uint64) (view domain.QueryView, err error)
}
