package core

import "crypto/md5"

type md5sum = [16]byte

/*
UnitCache

current cache is unlimited
*/
type UnitCache struct {
	// data md5: []Unit
	inner map[md5sum][]*Unit
}

func NewUnitCache() *UnitCache {
	return &UnitCache{
		make(map[md5sum][]*Unit),
	}
}

func (cache *UnitCache) Create(sum md5sum, value []*Unit) {
	// overwrite
	cache.inner[sum] = value
}

func (cache *UnitCache) CreateByData(data []byte, value []*Unit) {
	cache.Create(md5.Sum(data), value)
}

func (cache *UnitCache) Read(sum md5sum) []*Unit {
	ret, ok := cache.inner[sum]
	if !ok {
		return nil
	}
	return ret
}

func (cache *UnitCache) ReadByData(data []byte) []*Unit {
	return cache.Read(md5.Sum(data))
}
