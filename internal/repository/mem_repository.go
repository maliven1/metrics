package repository

type Storage struct {
	cache   Cache
	postgre Postgre
}

func NewCache(cache Cache, postgre Postgre) *Storage {
	return &Storage{cache: cache, postgre: postgre}
}

func (c *Storage) CheckCounter(key string) bool {
	return c.cache.CheckCounter(key)
}
func (c *Storage) CheckItemGauge(key string) bool {
	return c.cache.CheckItemGauge(key)
}

func (c *Storage) SetGauge(key string, value float64) {
	// Check if PostgreSQL is available
	if c.postgre != nil {
		if err := c.postgre.CheckConnection(); err == nil {
			// PostgreSQL is available, write to it
			c.postgre.SetGauge(key, value)

		}
	}
	// Fallback to memory storage
	c.cache.SetGauge(key, value)
}

func (c *Storage) SetCounter(key string, value int64) {
	// Check if PostgreSQL is available
	if c.postgre != nil {
		if err := c.postgre.CheckConnection(); err == nil {
			// PostgreSQL is available, write to it
			c.postgre.SetCounter(key, value)

		}
	}
	// Fallback to memory storage
	c.cache.SetCounter(key, value)
}

func (c *Storage) GetGauge() map[string]float64 {
	return c.cache.GetGauge()
}

func (c *Storage) GetItemGauge(s string) (string, float64) {
	return c.cache.GetItemGauge(s)
}
func (c *Storage) GetItemCounter(s string) (string, int64) {
	return c.cache.GetItemCounter(s)
}
func (c *Storage) GetCounter() map[string]int64 {
	return c.cache.GetCounter()
}

func (c *Storage) AddCounter(key string, value int64) bool {
	// Check if PostgreSQL is available
	if c.postgre != nil {
		if err := c.postgre.CheckConnection(); err == nil {
			// PostgreSQL is available, write to it
			c.postgre.AddCounter(key, value)
		}
	}
	// Fallback to memory storage
	return c.cache.AddCounter(key, value)
}
