package repository

import "context"

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
	SetGauge(key string, value float64, ctx context.Context) error
	SetCounter(key string, value int64, ctx context.Context) error
	GetAllGauges() (map[string]float64, error)
	GetAllCounters() (map[string]int64, error)
	GetItemGauge(key string) (string, float64, error)
	GetItemCounter(key string) (string, int64, error)
}

var usePostgre bool
