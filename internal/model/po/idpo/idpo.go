package idpo

import (
	"github.com/Me1onRind/go-demo/internal/model/po"
)

type IdCreator struct {
	po.BaseModel
	IdType IdType `gorm:"column:id_type"`
	Offset uint64 `gorm:"column:offset"`
	Step   uint32 `gorm:"column:step"`
}

func (i *IdCreator) TableName() string {
	return "id_creator_tab"
}
