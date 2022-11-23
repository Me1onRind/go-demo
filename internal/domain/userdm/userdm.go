package userdm

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	"github.com/Me1onRind/go-demo/internal/repository/userrepo"
)

type UserDomain struct {
	UserRepo *userrepo.UserRepo
}

func NewUserDomain() *UserDomain {
	return &UserDomain{
		UserRepo: userrepo.NewUserRepo(),
	}
}

func (u *UserDomain) CreateUser(ctx context.Context, user *userpo.User) (*userpo.User, error) {
	return nil, nil
}
