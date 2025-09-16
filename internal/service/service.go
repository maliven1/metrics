package service

import (
	"strconv"

	models "github.com/maliven1/metrics/internal/model"
)

type MemStorage interface {
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
	GetGauge() map[string]float64
	GetCounter() map[string]int64
	CheckCounter(key string) bool
	AddCounter(key string, value int64)
}

type Service struct {
	memStorage MemStorage
}

func NewService(m MemStorage) *Service {
	return &Service{memStorage: m}
}

func (s Service) CheckPath(pathSplit []string) int {
	if len(pathSplit) != 5 {
		return models.StatusNotFound
	}
	if float, err := strconv.ParseFloat(pathSplit[4], 64); pathSplit[2] == models.Gauge && err == nil && float != 0 {
		s.memStorage.SetGauge(pathSplit[3], float)
		return models.StatusOK
	} else if count, err := strconv.Atoi(pathSplit[4]); pathSplit[2] == models.Counter && err == nil && count != 0 {
		if s.memStorage.CheckCounter(pathSplit[3]) {
			s.memStorage.AddCounter(pathSplit[3], int64(count))
			return models.StatusOK
		}
		s.memStorage.SetCounter(pathSplit[3], int64(count))
		return models.StatusOK
	} else {
		return models.StatusBadRequest
	}
}
