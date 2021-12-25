package http_proto

type CreatePeridoicTaskReq struct {
	TaskName string `json:"task_name" binding:"required,min=1,max=128"`
	Cronspec string `json:"cronspec" binding:"required,min=1,max=64"`
	Status   uint8  `json:"status" binding:"required,oneof=0 1"`
}

type UpdatePeridoicTaskReq struct {
	ID       uint64 `json:"id" binding:"required,min=1"`
	TaskName string `json:"task_name" binding:"required,min=1,max=128"`
	Cronspec string `json:"cronspec" binding:"required,min=1,max=64"`
	Status   uint8  `json:"status" binding:"required,oneof=0 1"`
}

type GetPeridoicTaskByIDReq struct {
	ID uint64 `form:"id" binding:"required,min=1"`
}
