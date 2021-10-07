package kv_config_dao

const (
	KvConfigTablename = "kv_config_tab"
)

type KvConfigTab struct {
	ID        uint64 `gorm:"column:id" json:"id"`
	ConfigKey string `gorm:"column:config_key" json:"key"`
	Value     string `gorm:"column:value" json:"value"`
	ValueType uint8  `gorm:"column:value_type" json:"value_type"`
	Status    uint8  `gorm:"column:status" json:"status"`
	Ctime     uint32 `gorm:"autoCreateTime;column:ctime" json:"ctime"`
	Mtime     uint32 `gorm:"autoUpdateTime;column:mtime" json:"mtime"`
}

func (k *KvConfigTab) TableName() string {
	return KvConfigTablename
}
