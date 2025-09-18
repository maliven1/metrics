package repository

type Cache interface {
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
	GetGauge() map[string]float64
	GetCounter() map[string]int64
	CheckCounter(key string) bool
	AddCounter(key string, value int64)
	GetItemGauge(s string) (string, float64)
	GetItemCounter(s string) (string, int64)
}
type MemStorage struct {
	cache Cache
}

func NewCache(cache Cache) *MemStorage {
	return &MemStorage{cache: cache}
}

func (c *MemStorage) SetGauge(key string, value float64) {
	c.cache.SetGauge(key, value)
}

func (c *MemStorage) SetCounter(key string, value int64) {
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

func (c *MemStorage) CheckCounter(key string) bool {

	return c.cache.CheckCounter(key)
}

func (c *MemStorage) AddCounter(key string, value int64) {
	c.cache.AddCounter(key, value)
}
