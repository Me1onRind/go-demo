package respository

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
)

type ModelIfce interface {
	DBLabel() string
	TableName() string
	GetId() uint64
}

type BaseRepo[T ModelIfce] struct {
}

func (r *BaseRepo[T]) Create(ctx context.Context, value T) (T, error) {
	if err := mysql.GetWriteDB(ctx, value.DBLabel()).Create(value).Error; err != nil {
		logger.CtxErrorf(ctx, "Insert DB fail, value:[%+v], case:[%s]", value, err)
		var t T
		return t, gerror.WriteDBError.Wrap(err)
	}
	logger.CtxInfof(ctx, "Insert DB success, db:[%s], table:[%s], id:[%d]", value.DBLabel(), value.TableName(), value.GetId())
	return value, nil

}
