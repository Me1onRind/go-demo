package peridoic_task_protocol

type CreatePeridoicTaskReq struct {
	TaskName string `bind:"required,min=1,max=128"`
	Cronspec string `bind:"required,min=1,max=64"`
	Status   uint8  `bind:"required,oneof=0 1"`
}
