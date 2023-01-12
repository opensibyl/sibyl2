package core

import (
	"crypto/md5"
	"sync"
)

type md5sum = [16]byte

/*
UnitCache

current cache is unlimited
*/
type UnitCache struct {
	// data md5sum: []Unit
	inner sync.Map
}

func NewUnitCache() *UnitCache {
	return &UnitCache{}
}

func (cache *UnitCache) Create(sum md5sum, value []*Unit) {
	// overwrite
	cache.inner.Store(sum, value)
}

func (cache *UnitCache) CreateByData(data []byte, value []*Unit) {
	cache.Create(md5.Sum(data), value)
}

func (cache *UnitCache) Read(sum md5sum) []*Unit {
	ret, ok := cache.inner.Load(sum)
	if !ok {
		return nil
	}
	return ret.([]*Unit)
}

func (cache *UnitCache) ReadByData(data []byte) []*Unit {
	return cache.Read(md5.Sum(data))
}
