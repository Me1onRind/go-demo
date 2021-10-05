package store

import (
	"github.com/Me1onRind/go-demo/global/store/db_label"
	"gorm.io/gorm"
)

var (
	DBs = map[db_label.Label]*gorm.DB{}
)
