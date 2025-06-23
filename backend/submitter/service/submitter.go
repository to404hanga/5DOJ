package service

import (
	"5DOJ/pkg/constant/topic"
	"5DOJ/pkg/model"
	"5DOJ/submitter/domain"
	"5DOJ/submitter/global"
	"5DOJ/submitter/producer"
	"context"

	"github.com/google/uuid"
)

type SubmitterService struct {
	producer producer.Producer
}

var _ ISubmitterService = (*SubmitterService)(nil)

func NewSubmitterService(producer producer.Producer) *SubmitterService {
	return &SubmitterService{
		producer: producer,
	}
}

func (s *SubmitterService) Submit(ctx context.Context, contestId, problemId, userId uint64, lang, code, mode string) (recordId uint64, err error) {
	filenameWithoutExt := uuid.New().String()
	if err = model.InsertCode(ctx, model.Code{
		FilenameWithoutExt: filenameWithoutExt,
		Content:            code,
	}); err != nil {
		return
	}

	record := model.Record{
		UserId:    userId,
		ProblemId: problemId,
		ContestId: contestId,
		Language:  lang,
		CodeId:    filenameWithoutExt,
	}
	if err = global.MySQL.WithContext(ctx).Create(&record).Error; err != nil {
		return
	}

	send := topic.SubmitEvent{
		RecordId:           record.Id,
		ContestId:          contestId,
		ProblemId:          problemId,
		UserId:             userId,
		Language:           lang,
		FilenameWithoutExt: filenameWithoutExt,
		Code:               code,
		Mode:               mode,
	}
	err = s.producer.Produce(ctx, send)
	for err != nil {
		err = s.producer.Produce(ctx, send)
	}

	return record.Id, nil
}

func (s *SubmitterService) Query(ctx context.Context, recordId uint64) (view domain.QueryView, err error) {
	var record model.Record
	if err = global.MySQL.WithContext(ctx).Where("id = ?", recordId).First(&record).Error; err != nil {
		return
	}

	view = domain.QueryView{
		RecordId:      record.Id,
		ContestId:     record.ContestId,
		ProblemId:     record.ProblemId,
		UserId:        record.UserId,
		Language:      record.Language,
		Score:         record.Score,
		Result:        record.Result,
		TimeUsageMS:   record.TimeUsageMS,
		MemoryUsageKB: record.MemoryUsageKB,
		SubmitTime:    record.CreatedAt,
		UserName:      "", // TODO 后续接入用户服务
		ProblemTitle:  "", // TODO 后续接入题目服务
	}
	return
}
