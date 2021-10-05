package periodic_task

import (
	"github.com/Me1onRind/go-demo/internal/constant/localcache"
	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/Me1onRind/go-demo/internal/core/store"
	"github.com/Me1onRind/go-demo/internal/err_code"
)

type PeriodicTaskDao struct {
}

func NewPeriodicTaskDao() *PeriodicTaskDao {
	t := &PeriodicTaskDao{}
	return t
}

func (p *PeriodicTaskDao) LoadLocalCacheData() *common.Error {
	minID := uint64(0)
	limiter := 1000
	result := []*PeriodicTaskTab{}
	for {
		data := make([]*PeriodicTaskTab, 0, limiter)
		if err := store.ConfigDB.Where("id>?", minID).Limit(limiter).Order("id").Find(&data).Error; err != nil {
			return err_code.ReadDBError.WithErr(err)
		}
		if len(data) == 0 {
			break
		}
		result = append(result, data...)
		for _, v := range data {
			minID = v.ID
		}
	}

	if err := store.ConfigCache.Set(p.LocalCacheKey(), result); err != nil {
		return err_code.SetLocalCacheError.WithErr(err)
	}
	return nil
}

func (p *PeriodicTaskDao) LocalCacheKey() string {
	return localcache.PeriodicTask
}

func (p *PeriodicTaskDao) ListAllTask() []*PeriodicTaskTab {
	data, _ := store.ConfigCache.Get(localcache.PeriodicTask)
	if data != nil {
		tasks := data.([]*PeriodicTaskTab)
		return tasks
	}
	return nil
}
