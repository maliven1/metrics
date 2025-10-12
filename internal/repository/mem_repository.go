package repository

import "log"

type MemStorage struct {
	Storage
	cache Cache
}

func NewCache(cache Cache) *MemStorage {
	return &MemStorage{cache: cache}
}

func (c *MemStorage) CheckCounter(key string) bool {
	return c.cache.CheckCounter(key)
}
func (c *MemStorage) CheckItemGauge(key string) bool {
	return c.cache.CheckItemGauge(key)
}

func (c *MemStorage) SetGauge(key string, value float64) {
	if c.Storage.postgre != nil && c.Storage.postgre.CheckConnection() == nil {
		log.Println("22")
		c.Storage.postgre.SetGauge(key, value)
	}
	c.cache.SetGauge(key, value)
}

func (c *MemStorage) SetCounter(key string, value int64) {
	if c.Storage.postgre != nil && c.Storage.CheckConnection() == nil {
		log.Println("11")
		c.Storage.postgre.SetCounter(key, value)
	}
	c.cache.SetCounter(key, value)
}

func (c *MemStorage) GetGauge() map[string]float64 {
	return c.cache.GetGauge()
}

func (c *MemStorage) GetItemGauge(s string) (string, float64) {
	return c.cache.GetItemGauge(s)
}
func (c *MemStorage) GetItemCounter(s string) (string, int64) {
	return c.cache.GetItemCounter(s)
}
func (c *MemStorage) GetCounter() map[string]int64 {
	return c.cache.GetCounter()
}

func (c *MemStorage) AddCounter(key string, value int64) bool {

	if c.Storage.postgre != nil && c.Storage.CheckConnection() == nil {
		log.Println("33")
		c.Storage.postgre.SetCounter(key, value)
	}
	return c.cache.AddCounter(key, value)
}
