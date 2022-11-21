package userrepo

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/constant/dblabel"
	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
)

type UserRespo struct {
}

func NewUserRepo() *UserRespo {
	return &UserRespo{}
}

func (u *UserRespo) GetUserById(ctx context.Context, id uint64) (*userpo.User, error) {
	user := &userpo.User{}
	if err := mysql.GetDB(ctx, dblabel.Default).Take(user, "id=?", id).Error; err != nil {
		logger.CtxErrorf(ctx, "GetUserById fail, id:[%d], case:[%s]", id, err)
		return nil, gerror.ReadDBError.Wrap(err)
	}
	return user, nil
}

func (u *UserRespo) CreateUser(ctx context.Context, user *userpo.User) (*userpo.User, error) {
	if err := mysql.GetDB(ctx, dblabel.Default).Create(user).Error; err != nil {
		logger.CtxErrorf(ctx, "CreateUser fail, user:[%+v], case:[%s]", user, err)
		return nil, gerror.WriteDBError.Wrap(err)
	}
	logger.CtxInfof(ctx, "CreateUser success, id:[%d], name:[%s]", user.Id, user.Name)
	return user, nil
}
