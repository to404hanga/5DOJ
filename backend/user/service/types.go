package service

import (
	"5DOJ/user/domain"
	"context"
)

type IUserService interface {
	SignUp(ctx context.Context, userView domain.UserView) (err error)
	Login(ctx context.Context, uid, password string) (userView domain.UserView, err error)
	UpdatePassword(ctx context.Context, uid, password, confirmPassword string) (err error)
	GetInfoByUid(ctx context.Context, uid string) (userView domain.UserView, err error)
}
