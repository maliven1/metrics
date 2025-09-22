package storage

import (
	models "github.com/maliven1/metrics/internal/model"
)

type MemStorage struct {
	memCache *models.MemStorage
}

func NewMemStorage() *MemStorage {

	memStorage := &models.MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64)}

	return &MemStorage{
		memCache: memStorage,
	}
}

func (m MemStorage) SetGauge(key string, value float64) {
	m.memCache.Gauge[key] = value
}

func (m *MemStorage) SetCounter(key string, value int64) {
	m.memCache.Counter[key] = value

}

func (m *MemStorage) GetGauge() map[string]float64 {
	return m.memCache.Gauge
}
func (m *MemStorage) GetCounter() map[string]int64 {
	return m.memCache.Counter
}
func (m *MemStorage) GetItemCounter(s string) (string, int64) {
	v, ok := m.memCache.Counter[s]
	if ok {
		return s, v
	}
	return "", 0
}
func (m *MemStorage) GetItemGauge(s string) (string, float64) {
	v, ok := m.memCache.Gauge[s]
	if ok {
		return s, v
	}
	return "", 0
}

func (m *MemStorage) AddCounter(key string, value int64) bool {
	_, ok := m.memCache.Counter[key]
	if ok {
		m.memCache.Counter[key] += value
		return ok
	}

	return ok

}
