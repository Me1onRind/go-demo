package db_util

type ModelBase struct {
	Id    uint64 `json:"id" gorm:"column:id"`
	CTime uint32 `json:"ctime" gorm:"autoCreateTime;column:ctime"`
	MTime uint32 `json:"mtime" gorm:"autoUpdateTime;column:mtime"`
}
