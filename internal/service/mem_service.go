package service

import (
	"fmt"
	"strconv"

	models "github.com/maliven1/metrics/internal/model"
)

type MemService struct {
	memStorage MemRepo
}

func NewService(m MemRepo) *MemService {
	return &MemService{memStorage: m}
}

//go:generate mockgen -source=service.go -destination=mocks/mock.go
func (s MemService) AddStructMetric(metric models.Metrics) int {

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

func (s MemService) GetStructMetric(metric models.Metrics) (models.Metrics, int) {
	if metric.ID == "" {
		return metric, models.StatusBadRequest
	}

	if _, v := s.memStorage.GetItemGauge(metric.ID); metric.MType == models.Gauge && s.memStorage.CheckItemGauge(metric.ID) {
		metric.Value = &v

		return metric, models.StatusOK
	} else if _, v := s.memStorage.GetItemCounter(metric.ID); metric.MType == models.Counter && s.memStorage.CheckCounter(metric.ID) {

		metric.Delta = &v

		return metric, models.StatusOK
	}

	return metric, models.StatusNotFound
}

func (s MemService) CheckAddPath(pathSplit []string) int {
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

func (s MemService) GetMetric(pathSplit []string) (string, int) {
	if len(pathSplit) != 4 {
		return "", models.StatusNotFound
	}
	if name, v := s.memStorage.GetItemGauge(pathSplit[3]); pathSplit[2] == models.Gauge && name != "" {
		metrics := strconv.FormatFloat(v, 'f', -1, 64)
		return metrics, models.StatusOK
	} else if name, v := s.memStorage.GetItemCounter(pathSplit[3]); pathSplit[2] == models.Counter && name != "" {
		metrics := fmt.Sprint(v)
		return metrics, models.StatusOK
	}

	return "", models.StatusNotFound
}

func (s MemService) GetAllMetrics() (map[string]int64, map[string]float64) {

	counter := s.memStorage.GetCounter()
	gauge := s.memStorage.GetGauge()

	return counter, gauge
}
