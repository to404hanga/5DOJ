package grpc

import (
	problemv1 "5DOJ/api/proto/gen/problem/v1"
	"5DOJ/problem/domain"
	"5DOJ/problem/service"
	"context"

	"github.com/to404hanga/pkg404/gotools/transform"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProblemServiceServer struct {
	problemv1.UnimplementedProblemServiceServer
	svc service.IProblemService
}

func NewProblemServiceServer(svc service.IProblemService) *ProblemServiceServer {
	return &ProblemServiceServer{
		svc: svc,
	}
}

func (p *ProblemServiceServer) Register(srv grpc.ServiceRegistrar) {
	problemv1.RegisterProblemServiceServer(srv, p)
}

func (p *ProblemServiceServer) Get(ctx context.Context, req *problemv1.GetRequest) (*problemv1.GetResponse, error) {
	view, err := p.svc.Get(ctx, req.GetId())
	if err != nil {
		return &problemv1.GetResponse{}, err
	}
	return &problemv1.GetResponse{
		Problem: &problemv1.Problem{
			Id:            view.Id,
			Title:         view.Title,
			Level:         int32(view.Level),
			CreatedBy:     view.CreatedBy,
			UpdatedBy:     view.UpdatedBy,
			Enabled:       view.Enabled,
			TimeLimit:     int32(view.TimeLimit),
			MemoryLimit:   int32(view.MemoryLimit),
			TotalScore:    int32(view.TotalScore),
			TotalTestCase: int32(view.TotalTestCase),
			CreatedAt:     timestamppb.New(view.CreatedAt),
			UpdatedAt:     timestamppb.New(view.UpdatedAt),
			Markdown:      view.Markdown,
		},
	}, nil
}

func (p *ProblemServiceServer) GetTestCaseList(ctx context.Context, req *problemv1.GetTestCaseListRequest) (*problemv1.GetTestCaseListResponse, error) {
	view, err := p.svc.GetTestCaseList(ctx, req.GetId())
	if err != nil {
		return &problemv1.GetTestCaseListResponse{}, err
	}
	list := transform.SliceFromSlice[domain.TestCaseView, *problemv1.TestCase](view, func(i int, tcv domain.TestCaseView) *problemv1.TestCase {
		return &problemv1.TestCase{
			Id:        tcv.Id,
			Input:     tcv.Input,
			Output:    tcv.Output,
			Score:     int32(tcv.Score),
			CreatedBy: tcv.CreatedBy,
			UpdatedBy: tcv.UpdatedBy,
			Enabled:   tcv.Enabled,
		}
	})
	return &problemv1.GetTestCaseListResponse{
		Id:   req.GetId(),
		List: list,
	}, nil
}

func (p *ProblemServiceServer) GetList(ctx context.Context, req *problemv1.GetListRequest) (*problemv1.GetListResponse, error) {
	cursor, view, err := p.svc.GetList(ctx, int(req.GetSize()), req.GetCursor(), req.GetTitle())
	if err != nil {
		return &problemv1.GetListResponse{}, err
	}
	list := transform.SliceFromSlice[domain.ProblemView, *problemv1.Problem](view, func(i int, pv domain.ProblemView) *problemv1.Problem {
		return &problemv1.Problem{
			Id:            pv.Id,
			Title:         pv.Title,
			Level:         int32(pv.Level),
			CreatedBy:     pv.CreatedBy,
			UpdatedBy:     pv.UpdatedBy,
			Enabled:       pv.Enabled,
			TimeLimit:     int32(pv.TimeLimit),
			MemoryLimit:   int32(pv.MemoryLimit),
			TotalScore:    int32(pv.TotalScore),
			TotalTestCase: int32(pv.TotalTestCase),
			CreatedAt:     timestamppb.New(pv.CreatedAt),
			UpdatedAt:     timestamppb.New(pv.UpdatedAt),
			Markdown:      pv.Markdown,
		}
	})
	return &problemv1.GetListResponse{
		Size:   req.GetSize(),
		Cursor: cursor,
		List:   list,
	}, nil
}

func (p *ProblemServiceServer) Create(ctx context.Context, req *problemv1.CreateRequest) (*problemv1.CreateResponse, error) {
	pid, err := p.svc.Create(ctx, req.GetTitle(), int(req.GetLevel()), req.GetCreatedBy(), int(req.GetTimeLimit()), int(req.GetMemoryLimit()), req.GetMarkdown())
	if err != nil {
		return &problemv1.CreateResponse{}, err
	}
	return &problemv1.CreateResponse{
		Id: pid,
	}, nil
}

func (p *ProblemServiceServer) Update(ctx context.Context, req *problemv1.UpdateRequest) (*problemv1.UpdateResponse, error) {
	err := p.svc.Update(ctx, req.GetId(), req.GetTitle(), int(req.GetLevel()), req.GetUpdatedBy(), int(req.GetTimeLimit()), int(req.GetMemoryLimit()), req.GetMarkdown())
	return &problemv1.UpdateResponse{}, err
}

func (p *ProblemServiceServer) Enable(ctx context.Context, req *problemv1.EnableRequest) (*problemv1.EnableResponse, error) {
	err := p.svc.Enable(ctx, req.GetId(), req.GetUpdatedBy())
	return &problemv1.EnableResponse{}, err
}

func (p *ProblemServiceServer) Disable(ctx context.Context, req *problemv1.DisableRequest) (*problemv1.DisableResponse, error) {
	err := p.svc.Disable(ctx, req.GetId(), req.GetUpdatedBy())
	return &problemv1.DisableResponse{}, err
}

func (p *ProblemServiceServer) AppendTestCase(ctx context.Context, req *problemv1.AppendTestCaseRequest) (*problemv1.AppendTestCaseResponse, error) {
	tid, err := p.svc.AppendTestCase(ctx, req.GetId(), req.GetInput(), req.GetOutput(), int(req.GetScore()), req.GetCreatedBy())
	if err != nil {
		return &problemv1.AppendTestCaseResponse{}, err
	}
	return &problemv1.AppendTestCaseResponse{
		Pid: req.GetId(),
		Tid: tid,
	}, nil
}

func (p *ProblemServiceServer) UpdateTestCase(ctx context.Context, req *problemv1.UpdateTestCaseRequest) (*problemv1.UpdateTestCaseResponse, error) {
	err := p.svc.UpdateTestCase(ctx, req.GetPid(), req.GetTid(), req.GetInput(), req.GetOutput(), int(req.GetScore()), req.GetUpdatedBy())
	return &problemv1.UpdateTestCaseResponse{}, err
}

func (p *ProblemServiceServer) EnableTestCase(ctx context.Context, req *problemv1.EnableTestCaseRequest) (*problemv1.EnableTestCaseResponse, error) {
	err := p.svc.EnableTestCase(ctx, req.GetPid(), req.GetTid(), req.GetUpdatedBy())
	return &problemv1.EnableTestCaseResponse{}, err
}

func (p *ProblemServiceServer) DisableTestCase(ctx context.Context, req *problemv1.DisableTestCaseRequest) (*problemv1.DisableTestCaseResponse, error) {
	err := p.svc.DisableTestCase(ctx, req.GetPid(), req.GetTid(), req.GetUpdatedBy())
	return &problemv1.DisableTestCaseResponse{}, err
}
