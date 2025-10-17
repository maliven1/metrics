package service

import (
	"fmt"
	"strconv"

	models "github.com/maliven1/metrics/internal/model"
)

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

type MemService struct {
	memStorage MemRepo
}

func NewService(m MemRepo) *MemService {
	return &MemService{memStorage: m}
}

//go:generate mockgen -source=service.go -destination=mocks/mock.go
func (s MemService) AddStructMetric(metric models.Metrics) error {
	op := "service/AddStructMetric"
	if metric.MType == models.Gauge && metric.Value != nil {

		s.memStorage.SetGauge(metric.ID, *metric.Value)

	} else if metric.MType == models.Counter && metric.Delta != nil {

		if s.memStorage.AddCounter(metric.ID, *metric.Delta) {
			return nil
		}
		s.memStorage.SetCounter(metric.ID, *metric.Delta)

	} else {
		return fmt.Errorf("path:%s, err: BadRequest", op)
	}
	return nil
}

func (s MemService) GetStructMetric(metric models.Metrics) (models.Metrics, error) {
	op := "service/GetStructMetric"
	if metric.MType == models.Gauge && s.memStorage.CheckItemGauge(metric.ID) {
		_, v := s.memStorage.GetItemGauge(metric.ID)
		metric.Value = &v
		return metric, nil
	} else if metric.MType == models.Counter && s.memStorage.CheckCounter(metric.ID) {
		_, v := s.memStorage.GetItemCounter(metric.ID)
		metric.Delta = &v
		return metric, nil
	}
	return metric, fmt.Errorf("path%s, err: metric NotFound", op)
}

func (s MemService) CheckAddPath(pathSplit []string) error {
	op := "service/CheckAddPath"
	// Check if pathSplit has enough elements
	if len(pathSplit) < 5 {
		return fmt.Errorf("path: %s, err: BadRequest", op)
	}

	if float, err := strconv.ParseFloat(pathSplit[4], 64); pathSplit[2] == models.Gauge && err == nil {
		s.memStorage.SetGauge(pathSplit[3], float)
		return nil

	} else if count, err := strconv.Atoi(pathSplit[4]); pathSplit[2] == models.Counter && err == nil {
		if s.memStorage.AddCounter(pathSplit[3], int64(count)) {
			return nil
		}
		s.memStorage.SetCounter(pathSplit[3], int64(count))
		return nil
	}
	return fmt.Errorf("path: %s, err: BadRequest", op)
}

func (s MemService) GetMetric(pathSplit []string) (string, error) {
	op := "service/GetMetric"
	// Check if pathSplit has enough elements
	if len(pathSplit) < 4 {
		return "", fmt.Errorf("path: %s, err: BadRequest", op)
	}

	if pathSplit[2] == models.Gauge && s.memStorage.CheckItemGauge(pathSplit[3]) {
		_, v := s.memStorage.GetItemGauge(pathSplit[3])
		metrics := strconv.FormatFloat(v, 'f', -1, 64)
		return metrics, nil
	} else if pathSplit[2] == models.Counter && s.memStorage.CheckCounter(pathSplit[3]) {
		_, v := s.memStorage.GetItemCounter(pathSplit[3])
		metrics := fmt.Sprint(v)
		return metrics, nil
	}

	return "", fmt.Errorf("path: %s, err: metric NotFound", op)
}

func (s MemService) GetAllMetrics() (map[string]int64, map[string]float64) {

	counter := s.memStorage.GetCounter()
	gauge := s.memStorage.GetGauge()

	return counter, gauge
}

func (s MemService) SetMetrics(metrics []models.Metrics) {

	for _, v := range metrics {
		if v.MType == models.Gauge {
			s.memStorage.SetGauge(v.ID, *v.Value)

		} else if v.MType == models.Counter {
			if s.memStorage.AddCounter(v.ID, *v.Delta) {
				return
			}
			s.memStorage.SetCounter(v.ID, *v.Delta)

		}
	}

}
