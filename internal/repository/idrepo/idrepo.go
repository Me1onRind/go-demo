package idrepo

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/po/idpo"
	respository "github.com/Me1onRind/go-demo/internal/repository"
	"gorm.io/gorm"
)

type IdRepo struct {
	baseRepo *respository.BaseRepo[*idpo.IdCreator]
}

func NewIdRepo() *IdRepo {
	return &IdRepo{
		baseRepo: respository.NewBaseRepo[*idpo.IdCreator](),
	}
}

func (i *IdRepo) GetRecord(ctx context.Context, opts ...respository.Option) (*idpo.IdCreator, error) {
	return i.baseRepo.Take(ctx, opts...)
}

func WithIdType(idType idpo.IdType) respository.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id_type=?", idType)
	}
}

func (i *IdRepo) UpdateOffset(ctx context.Context, idType idpo.IdType, oldOffset uint64, step uint32) (int64, error) {
	model := &idpo.IdCreator{}
	db := mysql.GetWriteDB(ctx, model.DBLabel()).Model(model).Where("id_type=? AND offset=?", idType, oldOffset).Update("offset", oldOffset+uint64(step))
	if err := db.Error; err != nil {
		logger.CtxErrorf(ctx, "UpdateOffset fail, id_type:[%d], old_offset:[%d], step:[%d] case:[%s]", oldOffset, step, idType, err)
		return 0, gerror.WriteDBError.Wrap(err)
	}
	return db.RowsAffected, nil
}
