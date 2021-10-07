package periodic_task_dao

const (
	PeriodicTaskTableName = "periodic_task_tab"
)

type PeriodicTaskTab struct {
	ID       uint64 `json:"id" gorm:"column:id"`
	TaskName string `json:"task_name" gorm:"column:task_name;size=128"`
	Cronspec string `json:"cronspec" gorm:"column:cronspec;size=64"`
	Status   uint8  `json:"status" gorm:"column:status"`
	CTime    uint32 `json:"ctime" gorm:"autoCreateTime;column:ctime"`
	MTime    uint32 `json:"mtime" gorm:"autoUpdateTime;column:mtime"`
}

func (p *PeriodicTaskTab) TableName() string {
	return PeriodicTaskTableName
}

func (p *PeriodicTaskTab) GetID() uint64 {
	return p.ID
}
