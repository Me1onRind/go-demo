package userrepo

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (u *UserRepo) dbLabel() string {
	return gconfig.DynamicCfg.DefaultDB.GetLabel()
}

func (u *UserRepo) GetUserByUserId(ctx context.Context, userId uint64) (*userpo.User, error) {
	user := &userpo.User{}
	if err := mysql.GetReadDB(ctx, u.dbLabel()).Take(user, "user_id=?", userId).Error; err != nil {
		logger.CtxErrorf(ctx, "GetUserById fail, user_id:[%d], case:[%s]", userId, err)
		return nil, gerror.ReadDBError.Wrap(err)
	}
	return user, nil
}

func (u *UserRepo) CreateUser(ctx context.Context, user *userpo.User) (*userpo.User, error) {
	if err := mysql.GetWriteDB(ctx, u.dbLabel()).Create(user).Error; err != nil {
		logger.CtxErrorf(ctx, "CreateUser fail, user:[%+v], case:[%s]", user, err)
		return nil, gerror.WriteDBError.Wrap(err)
	}
	logger.CtxInfof(ctx, "CreateUser success, id:[%d], name:[%s]", user.Id, user.Name)
	return user, nil
}
