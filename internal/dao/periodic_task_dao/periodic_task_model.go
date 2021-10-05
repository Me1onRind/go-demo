package periodic_task_dao

var (
	PeriodicTaskTableName = "periodic_task_tab"
)

type PeriodicTaskTab struct {
	ID       uint64 `gorm:"column:id"`
	TaskName string `gorm:"column:task_name;size=128"`
	Cronspec string `gorm:"column:cronspec;size=64"`
	Status   uint8  `gorm:"column:status"`
	CTime    uint32 `gorm:"autoCreateTime;column:ctime"`
	MTime    uint32 `gorm:"autoUpdateTime;column:mtime"`
}

func (p *PeriodicTaskTab) TableName() string {
	return PeriodicTaskTableName
}

func (p *PeriodicTaskTab) GetID() uint64 {
	return p.ID
}
