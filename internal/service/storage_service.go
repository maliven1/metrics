// Package service
package service

import (
	"context"
)

type PostgreRepo interface {
	Close() error
	CheckConnection() error
	SetGauge(key string, value float64, ctx context.Context) error
	SetCounter(key string, value int64, ctx context.Context) error
	GetAllGauges() (map[string]float64, error)
	GetAllCounters() (map[string]int64, error)
}
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

func (s *PostgreService) CheckConnection() error {
	err := s.PostgreRepo.CheckConnection()
	if err != nil {
		return err
	}
	return nil
}
