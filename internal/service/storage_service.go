package service

import (
	"context"
	"net/http"

	models "github.com/maliven1/metrics/internal/model"
)

type PostgreService struct {
	PostgreRepo PostgreRepo
	MemRepo     MemRepo
}

func NewPostgreService(postgreRepo PostgreRepo, memRepo MemRepo) *PostgreService {
	return &PostgreService{
		PostgreRepo: postgreRepo,
		MemRepo:     memRepo,
	}
}

func (s *PostgreService) Close() {
	s.PostgreRepo.Close()

}

func (s *PostgreService) CheckConnection() int {
	err := s.PostgreRepo.CheckConnection()
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (s *PostgreService) SetMetrics(metrics []models.Metrics, ctx context.Context) (int, error) {
	if metrics == nil {
		return http.StatusBadRequest, nil
	}
	for _, v := range metrics {
		if v.MType == models.Gauge {
			err := s.PostgreRepo.SetGauge(v.ID, *v.Value, ctx)
			if err != nil {
				return http.StatusInternalServerError, err
			}

		} else if v.MType == models.Counter {
			err := s.PostgreRepo.SetCounter(v.ID, *v.Delta, ctx)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	gauges, err := s.PostgreRepo.GetAllGauges()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	for key, value := range gauges {
		s.MemRepo.SetGauge(key, value)

	}
	counters, err := s.PostgreRepo.GetAllCounters()
	if err != nil {
		return http.StatusInternalServerError, nil
	}
	for key, value := range counters {
		s.MemRepo.SetCounter(key, value)
	}
	return http.StatusOK, nil
}
