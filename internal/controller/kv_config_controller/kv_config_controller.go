package kv_config_controller

import (
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"github.com/Me1onRind/go-demo/internal/service/kv_config_service"
	protocl "github.com/Me1onRind/go-demo/protocol/http_proto/kv_config_protocol"
)

type KvConfigController struct {
	KvConfigService kv_config_service.KvConfigService
}

func NewKvConfigController() *KvConfigController {
	k := &KvConfigController{
		KvConfigService: *kv_config_service.NewKvConfigService(),
	}
	return k
}

func (k *KvConfigController) CreateKvConfig(ctx *ctm_context.Context, raw interface{}) (interface{}, *err_code.Error) {
	request := raw.(*protocl.CreateKVconfigReq)
	return k.KvConfigService.CreateKvConfig(ctx, request)
}

func (k *KvConfigController) UpdateKvConfig(ctx *ctm_context.Context, raw interface{}) (interface{}, *err_code.Error) {
	request := raw.(*protocl.UpdateKVconfigReq)
	return k.KvConfigService.UpdateKvConfig(ctx, request)
}

func (k *KvConfigController) GetKvConfigByID(ctx *ctm_context.Context, raw interface{}) (interface{}, *err_code.Error) {
	request := raw.(*protocl.GetKVconfigByIDReq)
	return k.KvConfigService.GetKvConfigByID(ctx, request.ID)
}
