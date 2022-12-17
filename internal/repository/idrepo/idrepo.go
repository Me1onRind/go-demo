package idrepo

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/po/idpo"
)

type IdRepo struct {
}

func NewIdRepo() *IdRepo {
	return &IdRepo{}
}

func (i *IdRepo) dbLabel() string {
	return gconfig.DynamicCfg.DefaultDB.GetLabel()
}

func (i *IdRepo) GetIdRecord(ctx context.Context, idType idpo.IdType) (*idpo.IdCreator, error) {
	var record idpo.IdCreator
	if err := mysql.GetReadDB(ctx, i.dbLabel()).Take(&record, "id_type=?", idType).Error; err != nil {
		logger.CtxErrorf(ctx, "GetIdRecord fail, id_type:[%d], case:[%s]", idType, err)
		return nil, gerror.ReadDBError.Wrap(err)
	}
	return &record, nil
}

func (i *IdRepo) UpdateOffset(ctx context.Context, idType idpo.IdType, oldOffset uint64, step uint32) (int64, error) {
	db := mysql.GetWriteDB(ctx, i.dbLabel()).Model(&idpo.IdCreator{}).Where("id_type=? AND offset=?", idType, oldOffset).Update("offset", oldOffset+uint64(step))
	if err := db.Error; err != nil {
		logger.CtxErrorf(ctx, "UpdateOffset fail, id_type:[%d], old_offset:[%d], step:[%d] case:[%s]", oldOffset, step, idType, err)
		return 0, gerror.WriteDBError.Wrap(err)
	}
	return db.RowsAffected, nil
}
