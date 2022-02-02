package kv_config_dao

import (
	"errors"
	"sync"

	"github.com/Me1onRind/go-demo/err_code"
	"github.com/Me1onRind/go-demo/infrastructure/ctm_context"
	"github.com/Me1onRind/go-demo/infrastructure/db_label"
	"github.com/Me1onRind/go-demo/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/constant"
	"github.com/Me1onRind/go-demo/internal/constant/kv_config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	k                  *KvConfigDao
	newKvConfigDaoOnce sync.Once
)

type KvConfigDao struct {
}

func NewKvConfigDao() *KvConfigDao {
	newKvConfigDaoOnce.Do(func() {
		k = &KvConfigDao{}
	})
	return k
}

func (k *KvConfigDao) GetKvConfigByID(ctx *ctm_context.Context, id uint64) (*KvConfigTab, *err_code.Error) {
	var kv KvConfigTab
	if err := ctx.GetDB(db_label.ConfigDB).Where("id=?", id).Take(&kv).Error; err != nil {
		return nil, err_code.ReadDBError.WithErr(err)
	}
	return &kv, nil
}

func (k *KvConfigDao) CreateKvConfig(ctx *ctm_context.Context, kv *KvConfigTab) *err_code.Error {
	if err := k.checkKvConfigValid(ctx, kv); err != nil {
		return err
	}
	var existKv KvConfigTab
	db := ctx.GetDB(db_label.ConfigDB)
	err := db.Where("config_key=?", kv.ConfigKey).Take(&existKv).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil {
			return err_code.DBRecordExistError.Withf("id:%d, config_key:%s", existKv.Id, existKv.ConfigKey)
		}
		return err_code.ReadDBError.WithErr(err)
	}
	if err := ctx.GetDB(db_label.ConfigDB).Create(kv).Error; err != nil {
		return err_code.WriteDBError.WithErr(err)
	}

	return nil
}

func (k *KvConfigDao) UpdateKvConfig(ctx *ctm_context.Context, kv *KvConfigTab) *err_code.Error {
	if err := k.checkKvConfigValid(ctx, kv); err != nil {
		return err
	}
	if err := ctx.GetDB(db_label.ConfigDB).Save(kv).Error; err != nil {
		return err_code.WriteDBError.WithErr(err)
	}

	return nil
}

func (k *KvConfigDao) LoadLocalcacheData(ctx *ctm_context.Context) (interface{}, *err_code.Error) {
	minID := uint64(0)
	limiter := 1000
	result := map[string]interface{}{}
	for {
		data := make([]*KvConfigTab, 0, limiter)
		if err := ctx.GetDB(db_label.ConfigDB).Where("id>?", minID).Limit(limiter).Order("id").Find(&data).Error; err != nil {
			return nil, err_code.ReadDBError.WithErr(err)
		}
		if len(data) == 0 {
			break
		}
		for _, v := range data {
			var value interface{}
			var ok bool
			if v.Status == constant.Enable {
				switch v.ValueType {
				case kv_config.Bool:
					value, ok = boolean(v.Value)
				case kv_config.Integer:
					value, ok = integer(v.Value)
				case kv_config.Float:
					value, ok = float(v.Value)
				case kv_config.List:
					value, ok = list(v.Value)
				case kv_config.Dict:
					value, ok = dict(v.Value)
				case kv_config.String:
					value = v.Value
					ok = true
				default:
					logger.CtxError(ctx, "Unknow value type", zap.Uint8("valueType", v.ValueType))
					continue
				}
				if !ok {
					logger.CtxError(ctx, "Can't convert", zap.Uint8("valueType", v.ValueType), zap.String("value", v.Value))
					continue
				}
				result[v.ConfigKey] = value
			}
			minID = v.Id
		}
	}
	return result, nil
}

//func (k *KvConfigDao) GetBool(ctx *ctm_context.Context, configKey string) (bool, *err_code.Error) {
//value, err := k.getCacheValue(ctx, configKey)
//if err != nil {
//return false, err
//}
//if realValue, ok := value.(bool); ok {
//return realValue, nil
//}
//return false, err_code.KvConfigKeyNotExistError.Withf("%v is not bool", value)
//}

//func (k *KvConfigDao) GetInt64(ctx *ctm_context.Context, configKey string) (int64, *err_code.Error) {
//value, err := k.getCacheValue(ctx, configKey)
//if err != nil {
//return 0, err
//}
//if realValue, ok := value.(int64); ok {
//return realValue, nil
//}
//return 0, err_code.KvConfigKeyNotExistError.Withf("%v is not int64", value)
//}

//func (k *KvConfigDao) GetFloat64(ctx *ctm_context.Context, configKey string) (float64, *err_code.Error) {
//value, err := k.getCacheValue(ctx, configKey)
//if err != nil {
//return 0, err
//}
//if realValue, ok := value.(float64); ok {
//return realValue, nil
//}
//return 0, err_code.KvConfigKeyNotExistError.Withf("%v is not float64", value)
//}

//func (k *KvConfigDao) GetString(ctx *ctm_context.Context, configKey string) (string, *err_code.Error) {
//value, err := k.getCacheValue(ctx, configKey)
//if err != nil {
//return "", err
//}
//if realValue, ok := value.(string); ok {
//return realValue, nil
//}
//return "", err_code.KvConfigKeyNotExistError.Withf("%v is not string", value)
//}

//func (k *KvConfigDao) GetDict(ctx *ctm_context.Context, configKey string) (map[string]interface{}, *err_code.Error) {
//value, err := k.getCacheValue(ctx, configKey)
//if err != nil {
//return nil, err
//}
//if realValue, ok := value.(map[string]interface{}); ok {
//return realValue, nil
//}
//return nil, err_code.KvConfigKeyNotExistError.Withf("%v is not dict", value)
//}

//func (k *KvConfigDao) GetList(ctx *ctm_context.Context, configKey string) ([]interface{}, *err_code.Error) {
//value, err := k.getCacheValue(ctx, configKey)
//if err != nil {
//return nil, err
//}

//if realValue, ok := value.([]interface{}); ok {
//return realValue, nil
//}
//return nil, err_code.KvConfigKeyNotExistError.Withf("%v is not list", value)

//}

func (k *KvConfigDao) checkKvConfigValid(ctx *ctm_context.Context, kv *KvConfigTab) *err_code.Error {
	switch kv.ValueType {
	case kv_config.Bool:
		if _, ok := boolean(kv.Value); !ok {
			return err_code.InvalidParamError.Withf("Can't convert %s to bool", kv.Value)
		}
	case kv_config.Integer:
		if _, ok := integer(kv.Value); !ok {
			return err_code.InvalidParamError.Withf("Can't convert %s to integer", kv.Value)
		}
	case kv_config.Float:
		if _, ok := float(kv.Value); !ok {
			return err_code.InvalidParamError.Withf("Can't convert %s to integer", kv.Value)
		}
	case kv_config.Dict:
		if _, ok := dict(kv.Value); !ok {
			return err_code.InvalidParamError.Withf("Can't convert %s to dict", kv.Value)
		}
	case kv_config.List:
		if _, ok := list(kv.Value); !ok {
			return err_code.InvalidParamError.Withf("Can't convert %s to list", kv.Value)
		}
	case kv_config.String:
	default:
		return err_code.InvalidParamError.Withf("Unknow value type:%d", kv.ValueType)
	}
	return nil
}
