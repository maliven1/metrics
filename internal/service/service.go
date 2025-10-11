package service

type PostgreRepo interface {
	Close() error
	CheckConnection() error
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
	AddCounter(key string, value int64) bool
}

type MemRepo interface {
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
	AddCounter(key string, value int64) bool
	GetItemGauge(s string) (string, float64)
	GetGauge() map[string]float64
	GetCounter() map[string]int64
	GetItemCounter(s string) (string, int64)
	CheckCounter(key string) bool
	CheckItemGauge(key string) bool
}
