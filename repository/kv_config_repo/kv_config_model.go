package kv_config_dao

import "github.com/Me1onRind/go-demo/infrastructure/util/db_util"

const (
	kvConfigTablename = "kv_config_tab"
)

type KvConfigTab struct {
	db_util.ModelBase
	ConfigKey string `gorm:"column:config_key" json:"key"`
	Value     string `gorm:"column:value" json:"value"`
	ValueType uint8  `gorm:"column:value_type" json:"value_type"`
	Status    uint8  `gorm:"column:status" json:"status"`
}

func (k *KvConfigTab) TableName() string {
	return kvConfigTablename
}
