package periodic_task_dao

import (
	"errors"

	"github.com/Me1onRind/go-demo/global/store"
	"github.com/Me1onRind/go-demo/global/store/db_label"
	"github.com/Me1onRind/go-demo/internal/constant/localcache"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"github.com/bluele/gcache"
	"gorm.io/gorm"
)

type PeriodicTaskDao struct {
}

func NewPeriodicTaskDao() *PeriodicTaskDao {
	t := &PeriodicTaskDao{}
	return t
}

func (p *PeriodicTaskDao) CreatePeriodicTask(ctx *ctm_context.Context, task *PeriodicTaskTab) *err_code.Error {
	var existTask PeriodicTaskTab
	db := ctx.DB(db_label.ConfigDB)
	err := db.Where("task_name=?", task.TaskName).Take(&existTask).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil {
			return err_code.DBRecordExistError.Withf("task_id:%d, task_name:%s", task.ID, task.TaskName)
		}
		return err_code.ReadDBError.WithErr(err)
	}

	if err := ctx.DB(db_label.ConfigDB).Create(task).Error; err != nil {
		return err_code.WriteDBError.WithErr(err)
	}
	return nil
}

func (p *PeriodicTaskDao) GetPeriodicTaskByID(ctx *ctm_context.Context, ID uint64) (*PeriodicTaskTab, *err_code.Error) {
	var task PeriodicTaskTab
	if err := ctx.DB(db_label.ConfigDB).Where("id=?", ID).Take(&task).Error; err != nil {
		return nil, err_code.ReadDBError.WithErr(err)
	}
	return &task, nil
}

func (p *PeriodicTaskDao) LoadLocalCacheData(ctx *ctm_context.Context) (interface{}, *err_code.Error) {
	minID := uint64(0)
	limiter := 1000
	result := []*PeriodicTaskTab{}
	for {
		data := make([]*PeriodicTaskTab, 0, limiter)
		if err := ctx.DB(db_label.ConfigDB).Where("id>?", minID).Limit(limiter).Order("id").Find(&data).Error; err != nil {
			return nil, err_code.ReadDBError.WithErr(err)
		}
		if len(data) == 0 {
			break
		}
		result = append(result, data...)
		for _, v := range data {
			minID = v.ID
		}
	}
	return result, nil
}

func (p *PeriodicTaskDao) LocalCacheKey() string {
	return localcache.PeriodicTask
}

func (p *PeriodicTaskDao) LocalCacheInstance() gcache.Cache {
	return store.ConfigCache
}

func (p *PeriodicTaskDao) ListAllTask() []*PeriodicTaskTab {
	data, _ := p.LocalCacheInstance().Get(localcache.PeriodicTask)
	if data != nil {
		tasks := data.([]*PeriodicTaskTab)
		return tasks
	}
	return nil
}
