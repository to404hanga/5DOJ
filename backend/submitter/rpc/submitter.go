package rpc

import (
	submitterv1 "5DOJ/api/proto/gen/submitter/v1"
	"5DOJ/submitter/service"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SubmitterServiceServer struct {
	submitterv1.UnimplementedSubmitterServiceServer
	svc service.ISubmitterService
}

func NewSubmitterServiceServer(svc service.ISubmitterService) *SubmitterServiceServer {
	return &SubmitterServiceServer{
		svc: svc,
	}
}

func (s *SubmitterServiceServer) Register(srv grpc.ServiceRegistrar) {
	submitterv1.RegisterSubmitterServiceServer(srv, s)
}

func (s *SubmitterServiceServer) Submit(ctx context.Context, req *submitterv1.SubmitRequest) (*submitterv1.SubmitResponse, error) {
	recordId, err := s.svc.Submit(ctx, req.GetContestId(), req.GetProblemId(), req.GetUserId(), req.GetLanguage(), req.GetCode(), int8(req.GetMode()))
	return &submitterv1.SubmitResponse{
		RecordId: recordId,
	}, err
}

func (s *SubmitterServiceServer) Query(ctx context.Context, req *submitterv1.QueryRequest) (*submitterv1.QueryResponse, error) {
	record, err := s.svc.Query(ctx, req.GetRecordId())
	return &submitterv1.QueryResponse{
		RecordId:      record.RecordId,
		ContestId:     record.ContestId,
		ProblemId:     record.ProblemId,
		UserId:        record.UserId,
		Language:      record.Language,
		Score:         int32(record.Score),
		Result:        record.Result,
		TimeUsageMS:   record.TimeUsageMS,
		MemoryUsageKB: record.MemoryUsageKB,
		Code:          record.Code,
		SubmitTime:    timestamppb.New(record.SubmitTime),
		UserName:      record.UserName,
	}, err
}
