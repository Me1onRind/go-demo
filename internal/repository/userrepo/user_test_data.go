package userrepo

import (
	"github.com/Me1onRind/go-demo/internal/model/po"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
)

var (
	UserJ = &userpo.User{
		BaseModel: po.BaseModel{
			Id:         1,
			CreateTime: 1680000000,
			UpdateTime: 1680000000,
		},
		UserId: 166,
		Email:  "test@google.com",
		Name:   "J",
	}
)
