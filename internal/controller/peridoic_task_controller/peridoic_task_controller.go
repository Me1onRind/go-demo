package peridoic_task_controller

import (
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"github.com/Me1onRind/go-demo/internal/service/peridoic_task_service"
	protocol "github.com/Me1onRind/go-demo/protocol/peridoic_task_protocol"
)

type PeridoicTaskController struct {
	PeridoicTaskService *peridoic_task_service.PeridoicTaskService
}

func NewPeridoicTaskController() *PeridoicTaskController {
	p := &PeridoicTaskController{
		PeridoicTaskService: peridoic_task_service.NewPeridoicTaskService(),
	}
	return p
}

func (p *PeridoicTaskController) CreatePeridoicTask(ctx *ctm_context.Context, raw interface{}) (interface{}, *err_code.Error) {
	request := raw.(*protocol.CreatePeridoicTaskReq)
	task, err := p.PeridoicTaskService.CreatePeridoicTask(ctx, request)
	if err != nil {
		return nil, err
	}

	return task, nil
}
