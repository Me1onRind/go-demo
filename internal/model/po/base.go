package po

type BaseModel struct {
	Id         uint64 `gorm:"column:id;primaryKey" json:"id"`
	CreateTime uint32 `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime uint32 `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}
