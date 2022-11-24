package useruc

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	"github.com/Me1onRind/go-demo/internal/repository/userrepo"
	"github.com/Me1onRind/go-demo/protocol/userproto"
)

type UserUsecase struct {
	UserRepo *userrepo.UserRepo
}

func NewUserUsecase() *UserUsecase {
	return &UserUsecase{
		UserRepo: userrepo.NewUserRepo(),
	}
}

func (u *UserUsecase) GetUserDetail(ctx context.Context, raw any) (any, error) {
	request := raw.(*userproto.GetUserDetailReq)
	return u.UserRepo.GetUserByUserId(ctx, request.UserId)
}

func (u *UserUsecase) CreateUser(ctx context.Context, raw any) (any, error) {
	request := raw.(*userproto.CreateUserReq)
	user := &userpo.User{
		Name: request.Name,
	}
	return u.UserRepo.CreateUser(ctx, user)
}
