package rpc

import (
	userv1 "5DOJ/api/proto/gen/user/v1"
	"5DOJ/user/domain"
	"5DOJ/user/service"
	"context"

	"google.golang.org/grpc"
)

type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
	svc service.IUserService
}

func NewUserServiceServer(svc service.IUserService) *UserServiceServer {
	return &UserServiceServer{
		svc: svc,
	}
}

func (u *UserServiceServer) Register(srv grpc.ServiceRegistrar) {
	userv1.RegisterUserServiceServer(srv, u)
}

func (u *UserServiceServer) SignUp(ctx context.Context, req *userv1.SignUpRequest) (*userv1.SignUpResponse, error) {
	return &userv1.SignUpResponse{}, u.svc.SignUp(ctx, domain.UserView{
		Uid:             req.GetUser().GetUid(),
		Name:            req.GetUser().GetName(),
		TelephoneNumber: req.GetUser().GetTelephoneNumber(),
		Gender:          req.GetUser().GetGender(),
	})
}

func (u *UserServiceServer) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	userView, err := u.svc.Login(ctx, req.GetUid(), req.GetPassword())
	if err != nil {
		return &userv1.LoginResponse{}, err
	}

	return &userv1.LoginResponse{
		User: &userv1.User{
			Uid:             userView.Uid,
			Name:            userView.Name,
			TelephoneNumber: userView.TelephoneNumber,
			Gender:          userView.Gender,
		},
	}, nil
}

func (u *UserServiceServer) UpdatePassword(ctx context.Context, req *userv1.UpdatePasswordRequest) (*userv1.UpdatePasswordResponse, error) {
	err := u.svc.UpdatePassword(ctx, req.GetUid(), req.GetPassword(), req.GetConfirmPassword())
	if err != nil {
		return &userv1.UpdatePasswordResponse{}, err
	}
	return &userv1.UpdatePasswordResponse{}, nil
}

func (u *UserServiceServer) GetInfoByUid(ctx context.Context, req *userv1.GetInfoByUidRequest) (*userv1.GetInfoByUidResponse, error) {
	userView, err := u.svc.GetInfoByUid(ctx, req.GetUid())
	if err != nil {
		return &userv1.GetInfoByUidResponse{}, err
	}
	return &userv1.GetInfoByUidResponse{
		User: &userv1.User{
			Uid:             userView.Uid,
			Name:            userView.Name,
			TelephoneNumber: userView.TelephoneNumber,
			Gender:          userView.Gender,
		},
	}, nil
}
