package storage

import (
	"sync"

	models "github.com/maliven1/metrics/internal/model"
)

type MemStorage struct {
	memCache *models.MemStorage
}

func NewMemStorage() *MemStorage {
	var m sync.RWMutex
	memStorage := &models.MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64), M: &m}

	return &MemStorage{
		memCache: memStorage,
	}
}

func (m MemStorage) SetGauge(key string, value float64) {
	m.memCache.M.Lock()
	defer m.memCache.M.Unlock()
	m.memCache.Gauge[key] = value
}

func (m *MemStorage) SetCounter(key string, value int64) {
	m.memCache.M.Lock()
	defer m.memCache.M.Unlock()
	m.memCache.Counter[key] = value

}

func (m *MemStorage) GetGauge() map[string]float64 {
	m.memCache.M.RLock()
	defer m.memCache.M.RUnlock()
	return m.memCache.Gauge
}
func (m *MemStorage) GetCounter() map[string]int64 {
	m.memCache.M.RLock()
	defer m.memCache.M.RUnlock()
	return m.memCache.Counter
}
func (m *MemStorage) GetItemCounter(s string) (string, int64) {
	m.memCache.M.RLock()
	defer m.memCache.M.RUnlock()
	v, ok := m.memCache.Counter[s]
	if ok {
		return s, v
	}
	return "", 0
}
func (m *MemStorage) GetItemGauge(key string) (string, float64) {
	m.memCache.M.RLock()
	defer m.memCache.M.RUnlock()
	v, ok := m.memCache.Gauge[key]
	if ok {
		return key, v
	}
	return "", 0
}
func (m *MemStorage) CheckItemGauge(key string) bool {
	m.memCache.M.RLock()
	defer m.memCache.M.RUnlock()
	_, ok := m.memCache.Gauge[key]
	if ok {
		return ok
	}
	return ok
}
func (m *MemStorage) CheckCounter(key string) bool {
	m.memCache.M.RLock()
	defer m.memCache.M.RUnlock()
	_, ok := m.memCache.Counter[key]
	if ok {
		return ok
	}

	return ok

}

func (m *MemStorage) AddCounter(key string, value int64) bool {
	m.memCache.M.Lock()
	defer m.memCache.M.Unlock()
	_, ok := m.memCache.Counter[key]
	if ok {
		m.memCache.Counter[key] += value
		return ok
	}

	return ok

}
