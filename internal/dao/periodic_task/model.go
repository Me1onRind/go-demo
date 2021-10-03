package periodic_task

var (
	PeriodicTaskTableName = "periodic_task_tab"
)

type PeriodicTaskTab struct {
	ID       uint64 `gorm:"column:id"`
	TaskName string `gorm:"column:task_name"`
	Cronspec string `gorm:"column:cronspec"`
	CTime    uint32 `gorm:"column:ctime"`
	MTime    uint32 `gorm:"column:mtime"`
}

func (p *PeriodicTaskTab) TableName() string {
	return PeriodicTaskTableName
}

func (p *PeriodicTaskTab) GetID() uint64 {
	return p.ID
}
