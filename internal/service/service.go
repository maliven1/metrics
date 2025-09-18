package service

import (
	"fmt"
	"strconv"

	models "github.com/maliven1/metrics/internal/model"
)

type MemStorage interface {
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
	CheckCounter(key string) bool
	AddCounter(key string, value int64)
	GetItemGauge(s string) (string, float64)
	GetGauge() map[string]float64
	GetCounter() map[string]int64
	GetItemCounter(s string) (string, int64)
}

type Service struct {
	memStorage MemStorage
}

func NewService(m MemStorage) *Service {
	return &Service{memStorage: m}
}

func (s Service) CheckAddPath(pathSplit []string) int {
	if len(pathSplit) != 5 {
		return models.StatusNotFound
	}
	if float, err := strconv.ParseFloat(pathSplit[4], 64); pathSplit[2] == models.Gauge && err == nil {
		s.memStorage.SetGauge(pathSplit[3], float)
		return models.StatusOK
	} else if count, err := strconv.Atoi(pathSplit[4]); pathSplit[2] == models.Counter && err == nil {
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
func (s Service) GetMetric(pathSplit []string) (string, int) {
	if len(pathSplit) != 4 {
		return "", models.StatusNotFound
	}
	if name, v := s.memStorage.GetItemGauge(pathSplit[3]); pathSplit[2] == models.Gauge && name != "" {
		metrics := fmt.Sprint(v)
		return metrics, models.StatusOK
	} else if name, v := s.memStorage.GetItemCounter(pathSplit[3]); pathSplit[2] == models.Counter && name != "" {
		metrics := fmt.Sprint(v)
		return metrics, models.StatusOK
	}

	return "", models.StatusNotFound
}

func (s Service) GetAllMetrics() (map[string]int64, map[string]float64) {

	counter := s.memStorage.GetCounter()
	gauge := s.memStorage.GetGauge()

	return counter, gauge
}
