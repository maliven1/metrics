package config

import models "github.com/maliven1/metrics/internal/model"

func InitMemStorage() models.MemStorage {
	MemStorage := models.MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64)}
	return MemStorage
}
