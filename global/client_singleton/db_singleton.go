package client_singleton

import (
	"github.com/Me1onRind/go-demo/infrastructure/db_label"
	"gorm.io/gorm"
)

var (
	DBs = map[db_label.Label]*gorm.DB{}
)
