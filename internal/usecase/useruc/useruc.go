package useruc

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/domain/userdm"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	"github.com/Me1onRind/go-demo/internal/repository/userrepo"
	"github.com/Me1onRind/go-demo/protocol/userproto"
)

type UserUsecase struct {
	UserDomain *userdm.UserDomain
	UserRepo   *userrepo.UserRepo
}

func NewUserUsecase() *UserUsecase {
	return &UserUsecase{
		UserRepo:   userrepo.NewUserRepo(),
		UserDomain: userdm.NewUserDomain(),
	}
}

func (u *UserUsecase) GetUserDetail(ctx context.Context, request *userproto.GetUserDetailReq) (any, error) {
	return u.UserRepo.GetUser(ctx, userrepo.WithUserId(request.UserId))
}

func (u *UserUsecase) CreateUser(ctx context.Context, request *userproto.CreateUserReq) (any, error) {
	user := &userpo.User{
		Name:  request.Name,
		Email: request.Email,
	}
	return u.UserDomain.CreateUser(ctx, user)
}
