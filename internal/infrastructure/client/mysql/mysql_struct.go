package mysql

import (
	"fmt"
	"sync"
	"sync/atomic"

	"gorm.io/gorm"
)

type gormSet struct {
	mutex *sync.Mutex
	infos map[string]*gormDBInfo
}

func newGormSet() *gormSet {
	return &gormSet{
		mutex: &sync.Mutex{},
		infos: make(map[string]*gormDBInfo),
	}
}

func (g *gormSet) register(dbLabel string, info *gormDBInfo) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.infos[dbLabel] = info
}

func (g *gormSet) getGormDBInfo(dbLabel string) *gormDBInfo {
	info := g.infos[dbLabel]
	if info == nil {
		panic(fmt.Sprintf("DB label:[%s] not exist", dbLabel))
	}
	return info
}

type gormDBInfo struct {
	writeDB     *gorm.DB
	readDBs     []*gorm.DB
	readDBIndex uint64
}

func newGormDBInfo(writeDB *gorm.DB, readDBs []*gorm.DB) *gormDBInfo {
	g := &gormDBInfo{
		writeDB: writeDB,
		readDBs: readDBs,
	}
	return g
}

func (g *gormDBInfo) getReadDB() *gorm.DB {
	if len(g.readDBs) > 0 {
		dbIndex := atomic.AddUint64(&g.readDBIndex, 1)
		return g.readDBs[dbIndex%uint64(len(g.readDBs))]
	}
	return g.writeDB
}

func (g *gormDBInfo) getWriteDB() *gorm.DB {
	return g.writeDB
}
