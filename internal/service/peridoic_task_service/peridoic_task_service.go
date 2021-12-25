package peridoic_task_service

import (
	"github.com/Me1onRind/go-demo/internal/dao/periodic_task_dao"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	protocol "github.com/Me1onRind/go-demo/protocol/http_proto/peridoic_task_protocol"
)

type PeridoicTaskService struct {
	PeridoicTaskDao *periodic_task_dao.PeriodicTaskDao
}

func NewPeridoicTaskService() *PeridoicTaskService {
	p := &PeridoicTaskService{
		PeridoicTaskDao: periodic_task_dao.NewPeriodicTaskDao(),
	}
	return p
}

func (p *PeridoicTaskService) CreatePeridoicTask(ctx *ctm_context.Context, req *protocol.CreatePeridoicTaskReq) (*periodic_task_dao.PeriodicTaskTab, *err_code.Error) {
	task := &periodic_task_dao.PeriodicTaskTab{
		TaskName: req.TaskName,
		Cronspec: req.Cronspec,
		Status:   req.Status,
	}

	if err := p.PeridoicTaskDao.CreatePeriodicTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (p *PeridoicTaskService) GetPeridoicTaskByID(ctx *ctm_context.Context, id uint64) (*periodic_task_dao.PeriodicTaskTab, *err_code.Error) {
	return p.PeridoicTaskDao.GetPeriodicTaskByID(ctx, id)
}

func (p *PeridoicTaskService) UpdatePeridoicTaskByID(ctx *ctm_context.Context, req *protocol.UpdatePeridoicTaskReq) (*periodic_task_dao.PeriodicTaskTab, *err_code.Error) {
	task, err := p.PeridoicTaskDao.GetPeriodicTaskByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	task.TaskName = req.TaskName
	task.Cronspec = req.Cronspec
	task.Status = req.Status

	if err := p.PeridoicTaskDao.UpdatePeriodicTask(ctx, task); err != nil {
		return nil, err
	}

	return task, err
}
