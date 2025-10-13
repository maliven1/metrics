package service

import models "github.com/maliven1/metrics/internal/model"

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

	gauges, err := s.PostgreRepo.GetAllGauges()
	if err != nil {
		return models.StatusInternalServerError
	}
	for key, value := range gauges {
		s.MemRepo.SetGauge(key, value)
	}
	counters, err := s.PostgreRepo.GetAllCounters()
	if err != nil {
		return models.StatusInternalServerError
	}
	for key, value := range counters {
		s.MemRepo.SetCounter(key, value)
	}
	return models.StatusOK
}
