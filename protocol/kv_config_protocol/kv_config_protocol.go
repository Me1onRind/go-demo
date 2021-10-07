package kv_conifg_protocol

type CreateKVconfigReq struct {
	ConfigKey string `json:"config_key" binding:"required,min=1,max=64"`
	Value     string `json:"value" binding:"required,min=1,max=2048"`
	ValueType uint8  `json:"value_type" binding:"required"`
	Status    uint8  `json:"status" binding:"required,oneof=0 1"`
}

type UpdateKVconfigReq struct {
	ID        uint64 `json:"id" binding:"required,min=1"`
	Value     string `json:"value" binding:"required,min=1,max=2048"`
	ValueType uint8  `json:"value_type" binding:"required"`
	Status    uint8  `json:"status" binding:"required,oneof=0 1"`
}

type GetKVconfigByIDReq struct {
	ID uint64 `form:"id" binding:"required,min=1"`
}
