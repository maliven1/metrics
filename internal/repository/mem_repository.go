package repository

import "context"

type MemStorage struct {
	postgre       Postgre
	cache         Cache
	usePostgreSQL bool
}

func NewCache(cache Cache, usePostgreSQL bool, postgre Postgre) *MemStorage {

	return &MemStorage{cache: cache, usePostgreSQL: usePostgreSQL, postgre: postgre}
}

func (c *MemStorage) CheckCounter(key string) bool {
	return c.cache.CheckCounter(key)
}
func (c *MemStorage) CheckItemGauge(key string) bool {
	return c.cache.CheckItemGauge(key)
}

func (c *MemStorage) SetGauge(key string, value float64) {
	if c.usePostgreSQL {
		c.postgre.SetGauge(key, value, context.Background())
	}
	c.cache.SetGauge(key, value)
}

func (c *MemStorage) SetCounter(key string, value int64) {
	if c.usePostgreSQL {
		c.postgre.SetCounter(key, value, context.Background())
	}
	c.cache.SetCounter(key, value)
}

func (c *MemStorage) GetGauge() map[string]float64 {
	if c.usePostgreSQL {
		m, _ := c.postgre.GetAllGauges()
		return m
	}
	return c.cache.GetGauge()
}

func (c *MemStorage) GetItemGauge(s string) (string, float64) {
	if c.usePostgreSQL {
		k, v, _ := c.postgre.GetItemGauge(s)
		return k, v
	}
	return c.cache.GetItemGauge(s)
}
func (c *MemStorage) GetItemCounter(s string) (string, int64) {
	if c.usePostgreSQL {
		k, v, _ := c.postgre.GetItemCounter(s)
		return k, v
	}
	return c.cache.GetItemCounter(s)
}
func (c *MemStorage) GetCounter() map[string]int64 {
	if c.usePostgreSQL {
		m, _ := c.postgre.GetAllCounters()
		return m
	}
	return c.cache.GetCounter()
}

func (c *MemStorage) AddCounter(key string, value int64) bool {
	if c.usePostgreSQL {
		c.postgre.SetCounter(key, value, context.Background())
	}

	return c.cache.AddCounter(key, value)
}
