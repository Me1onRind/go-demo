package periodic_task_dao

import (
	"errors"
	"sync"

	"github.com/Me1onRind/go-demo/err_code"
	"github.com/Me1onRind/go-demo/infrastructure/ctm_context"
	"github.com/Me1onRind/go-demo/infrastructure/db_label"
	"gorm.io/gorm"
)

var (
	newPeridoicTaskDaoOnce sync.Once
	t                      *PeriodicTaskDao
)

type PeriodicTaskDao struct {
}

func NewPeriodicTaskDao() *PeriodicTaskDao {
	newPeridoicTaskDaoOnce.Do(func() {
		t = &PeriodicTaskDao{}
	})
	return t
}

func (p *PeriodicTaskDao) CreatePeriodicTask(ctx *ctm_context.Context, task *PeriodicTaskTab) *err_code.Error {
	if err := p.checkPeriodicTaskNoExist(ctx, task.TaskName, task.Id); err != nil {
		return err
	}

	if err := ctx.GetDB(db_label.ConfigDB).Create(task).Error; err != nil {
		return err_code.WriteDBError.WithErr(err)
	}

	return nil
}

func (p *PeriodicTaskDao) GetPeriodicTaskByID(ctx *ctm_context.Context, id uint64) (*PeriodicTaskTab, *err_code.Error) {
	var task PeriodicTaskTab
	if err := ctx.GetDB(db_label.ConfigDB).Where("id=?", id).Take(&task).Error; err != nil {
		return nil, err_code.ReadDBError.WithErr(err)
	}
	return &task, nil
}

func (p *PeriodicTaskDao) UpdatePeriodicTask(ctx *ctm_context.Context, task *PeriodicTaskTab) *err_code.Error {
	if err := p.checkPeriodicTaskNoExist(ctx, task.TaskName, task.Id); err != nil {
		return err
	}

	if err := ctx.GetDB(db_label.ConfigDB).Save(task).Error; err != nil {
		return err_code.WriteDBError.WithErr(err)
	}

	return nil
}

func (p *PeriodicTaskDao) checkPeriodicTaskNoExist(ctx *ctm_context.Context, taskName string, skipID uint64) *err_code.Error {
	var existTask PeriodicTaskTab
	db := ctx.GetDB(db_label.ConfigDB)
	err := db.Where("task_name=?", taskName).Where("id!=?", skipID).Take(&existTask).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil {
			return err_code.DBRecordExistError.Withf("task_id:%d, task_name:%s", existTask.Id, existTask.TaskName)
		}
		return err_code.ReadDBError.WithErr(err)
	}
	return nil
}

func (p *PeriodicTaskDao) LoadLocalcacheData(ctx *ctm_context.Context) (interface{}, *err_code.Error) {
	minID := uint64(0)
	limiter := 1000
	result := []*PeriodicTaskTab{}
	for {
		data := make([]*PeriodicTaskTab, 0, limiter)
		if err := ctx.GetDB(db_label.ConfigDB).Where("id>?", minID).Limit(limiter).Order("id").Find(&data).Error; err != nil {
			return nil, err_code.ReadDBError.WithErr(err)
		}
		if len(data) == 0 {
			break
		}
		result = append(result, data...)
		for _, v := range data {
			minID = v.Id
		}
	}
	return result, nil
}

//func (p *PeriodicTaskDao) ListAllTask() []*PeriodicTaskTab {
//data, _ := p.Localcache.Get(localcache.PeriodicTask)
//if data != nil {
//tasks := data.([]*PeriodicTaskTab)
//return tasks
//}
//return nil
//}
