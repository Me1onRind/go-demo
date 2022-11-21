package po

type BaseModel struct {
	Id         uint64 `gorm:"column:id;primaryKey"`
	CreateTime uint32 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime uint32 `gorm:"column:update_time;autoUpdateTime"`
}
