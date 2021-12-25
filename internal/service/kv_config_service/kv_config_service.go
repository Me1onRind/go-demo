package kv_config_service

import (
	"github.com/Me1onRind/go-demo/internal/dao/kv_config_dao"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	protocol "github.com/Me1onRind/go-demo/protocol/http_proto/kv_config_protocol"
)

type KvConfigService struct {
	KvConfigDao *kv_config_dao.KvConfigDao
}

func NewKvConfigService() *KvConfigService {
	k := &KvConfigService{
		KvConfigDao: kv_config_dao.NewKvConfigDao(),
	}
	return k
}

func (k *KvConfigService) CreateKvConfig(ctx *ctm_context.Context, req *protocol.CreateKVconfigReq) (*kv_config_dao.KvConfigTab, *err_code.Error) {
	kv := &kv_config_dao.KvConfigTab{
		ConfigKey: req.ConfigKey,
		Value:     req.Value,
		ValueType: req.ValueType,
		Status:    req.Status,
	}
	if err := k.KvConfigDao.CreateKvConfig(ctx, kv); err != nil {
		return nil, err
	}
	return kv, nil
}

func (k *KvConfigService) UpdateKvConfig(ctx *ctm_context.Context, req *protocol.UpdateKVconfigReq) (*kv_config_dao.KvConfigTab, *err_code.Error) {
	kv, err := k.KvConfigDao.GetKvConfigByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	kv.Value = req.Value
	kv.ValueType = req.ValueType
	kv.Status = req.Status

	if err := k.KvConfigDao.UpdateKvConfig(ctx, kv); err != nil {
		return nil, err
	}
	return kv, nil
}

func (k *KvConfigService) GetKvConfigByID(ctx *ctm_context.Context, id uint64) (*kv_config_dao.KvConfigTab, *err_code.Error) {
	return k.KvConfigDao.GetKvConfigByID(ctx, id)
}
