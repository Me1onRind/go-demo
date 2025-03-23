package userpo

import (
	"github.com/Me1onRind/go-demo/internal/model/po"
	"github.com/Me1onRind/go-demo/protocol/userproto"
)

type User struct {
	po.BaseModel
	UserId uint64 `gorm:"column:user_id" json:"user_id"`
	Email  string `gorm:"column:email"`
	Name   string `gorm:"column:name" json:"name"`
}

func (u *User) TableName() string {
	return "user_tab"
}

func (u *User) ToProtocoUser() *userproto.UserInfo {
	return &userproto.UserInfo{
		UserId: u.UserId,
		Name:   u.Name,
		Email:  u.Email,
	}
}
