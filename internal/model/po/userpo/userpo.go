package userpo

import "github.com/Me1onRind/go-demo/internal/model/po"

type User struct {
	po.BaseModel
	UserId uint64 `gorm:"column:user_id" json:"user_id"`
	Name   string `gorm:"column:name" json:"name"`
}

func (u *User) TableName() string {
	return "user_tab"
}
