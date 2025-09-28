package service

import (
	"fmt"
	"strconv"

	models "github.com/maliven1/metrics/internal/model"
)

type MemStorage interface {
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
	AddCounter(key string, value int64) bool
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

//go:generate mockgen -source=service.go -destination=mocks/mock.go
func (s Service) AddStructMetric(metric models.Metrics) int {

	if metric.MType == models.Gauge && metric.Value != nil {
		s.memStorage.SetGauge(metric.ID, *metric.Value)
		return models.StatusOK
	} else if metric.MType == models.Counter && metric.Delta != nil {
		if s.memStorage.AddCounter(metric.ID, *metric.Delta) {
			return models.StatusOK
		}
		s.memStorage.SetCounter(metric.ID, *metric.Delta)
		return models.StatusOK
	} else {
		return models.StatusBadRequest
	}
}

func (s Service) GetStructMetric(metric models.Metrics) (models.Metrics, int) {

	if name, v := s.memStorage.GetItemGauge(metric.ID); metric.MType == models.Gauge && name != "" {
		metric.Value = &v

		return metric, models.StatusOK
	} else if name, v := s.memStorage.GetItemCounter(metric.ID); metric.MType == models.Counter && name != "" {

		metric.Delta = &v

		return metric, models.StatusOK
	}

	return metric, models.StatusNotFound
}

func (s Service) CheckAddPath(pathSplit []string) int {
	if len(pathSplit) != 5 {
		return models.StatusNotFound
	}
	if float, err := strconv.ParseFloat(pathSplit[4], 64); pathSplit[2] == models.Gauge && err == nil {
		s.memStorage.SetGauge(pathSplit[3], float)
		return models.StatusOK
	} else if count, err := strconv.Atoi(pathSplit[4]); pathSplit[2] == models.Counter && err == nil {
		if s.memStorage.AddCounter(pathSplit[3], int64(count)) {
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
