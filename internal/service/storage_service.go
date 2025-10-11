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
