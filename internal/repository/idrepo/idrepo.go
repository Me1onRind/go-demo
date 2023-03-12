package idrepo

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/po/idpo"
	"gorm.io/gorm"
)

type IdRepo struct {
}

func NewIdRepo() *IdRepo {
	return &IdRepo{}
}

func (i *IdRepo) DBLabel() string {
	return gconfig.DynamicCfg.DefaultDB.GetLabel()
}

func (i *IdRepo) GetRecord(ctx context.Context, opts ...mysql.Option) (*idpo.IdCreator, error) {
	return mysql.Take[idpo.IdCreator](ctx, i.DBLabel(), opts...)
}

func WithIdType(idType idpo.IdType) mysql.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id_type=?", idType)
	}
}

func (i *IdRepo) UpdateOffset(ctx context.Context, idType idpo.IdType, oldOffset uint64, step uint32) (int64, error) {
	model := &idpo.IdCreator{}
	db := mysql.GetWriteDB(ctx, i.DBLabel()).Model(model).Where("id_type=? AND offset=?", idType, oldOffset).Update("offset", oldOffset+uint64(step))
	if err := db.Error; err != nil {
		logger.CtxErrorf(ctx, "UpdateOffset fail, id_type:[%d], old_offset:[%d], step:[%d] case:[%s]", oldOffset, step, idType, err)
		return 0, gerror.WriteDBError.Wrap(err)
	}
	return db.RowsAffected, nil
}
