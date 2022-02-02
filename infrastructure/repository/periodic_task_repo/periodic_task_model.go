package periodic_task_dao

import "github.com/Me1onRind/go-demo/infrastructure/util/db_util"

const (
	periodicTaskTableName = "periodic_task_tab"
)

type PeriodicTaskTab struct {
	db_util.ModelBase
	TaskName string `json:"task_name" gorm:"column:task_name;size=128"`
	Cronspec string `json:"cronspec" gorm:"column:cronspec;size=64"`
	Status   uint8  `json:"status" gorm:"column:status"`
}

func (p *PeriodicTaskTab) TableName() string {
	return periodicTaskTableName
}
