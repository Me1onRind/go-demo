package iddm

import (
	"context"
	"errors"
	"sync"

	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/po/idpo"
	"github.com/Me1onRind/go-demo/internal/repository/idrepo"
)

var (
	once                 = sync.Once{}
	instance             *idDomain
	ErrGetPoolFromDBFail = errors.New("Get id pool from database fail, may try too much times")
	ErrPoolStepIsZero    = errors.New("Id pool step is zero")
)

type idDomain struct {
	IdRepo *idrepo.IdRepo
	Mutex  *sync.RWMutex
	pools  map[idpo.IdType]*idPool
}

type idPool struct {
	Mutex  *sync.Mutex
	Max    uint64
	Offset uint64
	Count  uint64
}

func (i *idPool) getId() (uint64, bool) {
	id := i.Count + i.Offset
	if id >= i.Max {
		return 0, false
	}
	i.Count++
	return id, true
}

type IdDomain interface {
	GetId(ctx context.Context, idType idpo.IdType, maxRetry int) (uint64, error)
}

func NewIdDomain() IdDomain {
	once.Do(func() {
		instance = &idDomain{
			pools:  map[idpo.IdType]*idPool{},
			Mutex:  &sync.RWMutex{},
			IdRepo: idrepo.NewIdRepo(),
		}
	})
	return instance
}

func (i *idDomain) getPool(idType idpo.IdType) *idPool {
	i.Mutex.RLock()
	pool := i.pools[idType]
	i.Mutex.RUnlock()
	if pool != nil {
		return pool
	}

	pool = &idPool{
		Mutex: &sync.Mutex{},
	}
	i.Mutex.Lock()
	i.pools[idType] = pool
	i.Mutex.Unlock()
	return pool
}

func (i *idDomain) GetId(ctx context.Context, idType idpo.IdType, maxTry int) (uint64, error) {
	if maxTry <= 0 {
		maxTry = 1
	}

	pool := i.getPool(idType)
	pool.Mutex.Lock()
	defer pool.Mutex.Unlock()

	id, ok := pool.getId()
	if ok {
		logger.CtxDebugf(ctx, "generate idtype:%s, id:%d", idType, id)
		return id, nil
	}

	logger.CtxDebugf(ctx, "idtype:%s id used up, try pull", idType)
	for j := range maxTry {
		record, err := i.IdRepo.GetRecord(ctx, idrepo.WithIdType(idType))
		if err != nil {
			return 0, gerror.GenerateIdError.Infof("id_type:%s", idType).Wrap(err)
		}

		if record.Step == 0 {
			logger.CtxErrorf(ctx, "Id pool is invalid, idType:[%d]", idType)
			return 0, gerror.GenerateIdError.Infof("id_type:%s", idType).Wrap(err)
		}

		rows, err := i.IdRepo.UpdateOffset(ctx, idType, record.Offset, record.Step)
		if err != nil {
			return 0, gerror.GenerateIdError.Wrap(ErrPoolStepIsZero)
		}

		if rows == 0 {
			logger.CtxWarnf(ctx, "UpdateOffset affectRows is zero, id_type:[%d], Offset[%d], step:[%d], times[%d], retryTimes[%d]",
				idType, record.Offset, record.Step, j+1, maxTry)
			continue
		}

		pool.Max = record.Offset + uint64(record.Step)
		pool.Count = 0
		pool.Offset = record.Offset
		id, _ := pool.getId()
		return id, nil
	}

	return 0, gerror.GenerateIdError.Infof("id_type:%s", idType).Wrap(ErrGetPoolFromDBFail)
}
