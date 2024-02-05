package ants

import (
	"sync"

	"github.com/panjf2000/ants/v2"
)

// Ants 对应的结构体
var Ants = &antsStruct{}
var antMutex = &sync.RWMutex{}

type antsStruct struct {
	Pools map[string]*ants.Pool
}

// 提交函数
func (p *antsStruct) Submit(poolName string, task func(), poolSize int) error {
	var pool *ants.Pool
	if poolSize == 0 {
		poolSize = 30
	}
	if p.Pools == nil {
		p.Pools = make(map[string]*ants.Pool)
	}
	antMutex.RLock()
	_, ok := p.Pools[poolName]
	antMutex.RUnlock()
	if !ok {
		options := ants.WithOptions(ants.Options{
			PreAlloc: true,
		})
		pn, _ := ants.NewPool(poolSize, options)
		antMutex.Lock()
		p.Pools[poolName] = pn
		antMutex.Unlock()
	}
	antMutex.RLock()
	pool = p.Pools[poolName]
	antMutex.RUnlock()
	err := pool.Submit(task)

	return err
}
