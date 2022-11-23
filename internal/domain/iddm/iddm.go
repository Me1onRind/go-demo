package iddm

import (
	"context"
	"sync"

	"github.com/Me1onRind/go-demo/internal/model/po/idpo"
	"github.com/Me1onRind/go-demo/internal/repository/idrepo"
)

var once = sync.Once{}
var instance *IdDomain

type IdDomain struct {
	IdRepo *idrepo.IdRepo
	Step   uint64
	Max    uint64
	Offset uint64
	Count  uint64
	Mutex  *sync.Mutex
}

func NewIdDomain() *IdDomain {
	once.Do(func() {
		instance = &IdDomain{
			Mutex: &sync.Mutex{},
		}
	})
	return instance
}

func (i *IdDomain) GetId(ctx context.Context, idType idpo.IdType, maxRetry int) (uint64, error) {
	var (
		err error
		id  uint64
	)
	i.Mutex.Lock()
	defer i.Mutex.Unlock()
	for j := 0; j < maxRetry; j++ {
		tmpId := i.Count + i.Offset
		if tmpId < i.Max {
			id = tmpId
		}
	}
	return id, err
}
