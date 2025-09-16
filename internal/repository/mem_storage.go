package repository

import (
	"github.com/maliven1/metrics/internal/config"
	models "github.com/maliven1/metrics/internal/model"
)

type MemStorage struct {
	memCache *models.MemStorage
}

type Storage interface {
	MemStorage
}

func InitCache() *MemStorage {
	memStorage := config.InitMemStorage()

	return &MemStorage{
		memCache: &memStorage,
	}
}

func (c *MemStorage) SetGauge(key string, value float64) {
	c.memCache.Gauge[key] = value

}

func (c *MemStorage) SetCounter(key string, value int64) {
	c.memCache.Counter[key] = value

}

func (c *MemStorage) GetGauge() map[string]float64 {
	return c.memCache.Gauge
}
func (c *MemStorage) GetCounter() map[string]int64 {
	return c.memCache.Counter
}
func (c *MemStorage) CheckCounter(key string) bool {
	_, ok := c.memCache.Counter[key]
	if ok {
		return ok
	}

	return ok
}

func (c *MemStorage) AddCounter(key string, value int64) {
	c.memCache.Counter[key] = +value

}
