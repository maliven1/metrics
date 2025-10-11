package repository

type Cache interface {
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
	GetGauge() map[string]float64
	GetCounter() map[string]int64
	AddCounter(key string, value int64) bool
	GetItemGauge(key string) (string, float64)
	GetItemCounter(key string) (string, int64)
	CheckCounter(key string) bool
	CheckItemGauge(key string) bool
}
type Postgre interface {
	Close() error
	CheckConnection() error
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
	AddCounter(key string, value int64) bool
}
