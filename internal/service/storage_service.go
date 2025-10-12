package service

import models "github.com/maliven1/metrics/internal/model"

type PostgreService struct {
	PostgreRepo PostgreRepo
}

func NewPostgreService(repo PostgreRepo) *PostgreService {
	return &PostgreService{PostgreRepo: repo}
}

func (s *PostgreService) Close() {
	s.PostgreRepo.Close()

}

func (s *PostgreService) CheckConnection() int {
	err := s.PostgreRepo.CheckConnection()
	if err != nil {
		return models.StatusInternalServerError
	}
	return models.StatusOK
}

func (s *PostgreService) SetMetrics(metrics []models.Metrics) int {
	if metrics == nil {
		return models.StatusBadRequest
	}
	for _, v := range metrics {
		if v.MType == models.Gauge {
			s.PostgreRepo.SetGauge(v.ID, *v.Value)
		} else if v.MType == models.Counter {
			s.PostgreRepo.SetCounter(v.ID, *v.Delta)
		}
	}
	return models.StatusOK
}
