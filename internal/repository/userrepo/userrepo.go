package userrepo

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	respository "github.com/Me1onRind/go-demo/internal/repository"
)

type UserRepo struct {
	*respository.BaseRepo[*userpo.User]
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (*userpo.User, error) {
	user := &userpo.User{}
	if err := mysql.GetReadDB(ctx, user.DBLabel()).Take(user, "email=?", email).Error; err != nil {
		logger.CtxErrorf(ctx, "GetUserByEmail fail, email:[%s], case:[%s]", email, err)
		return nil, gerror.ReadDBError.Wrap(err)
	}
	return user, nil
}

func (u *UserRepo) GetUserByUserId(ctx context.Context, userId uint64) (*userpo.User, error) {
	user := &userpo.User{}
	if err := mysql.GetReadDB(ctx, user.DBLabel()).Take(user, "user_id=?", userId).Error; err != nil {
		logger.CtxErrorf(ctx, "GetUserById fail, user_id:[%d], case:[%s]", userId, err)
		return nil, gerror.ReadDBError.Wrap(err)
	}
	return user, nil
}
