package repository

import (
	"context"
	"time"

	"github.com/maliven1/metrics/internal/repository/pgerrors"
)

type Storage struct {
	postgre Postgre
}

func NewStorage(postgre Postgre) *Storage {
	return &Storage{postgre: postgre}
}

func LinearBackoff(callback func() error) error {
	const maxRetries = 3
	var lastErr error
	const interval = 2
	LinearBackoff := 1
	classifier := pgerrors.NewPostgresErrorClassifier()
	for i := 0; i <= maxRetries; i++ {
		time.Sleep(time.Duration(LinearBackoff) * time.Second)
		LinearBackoff += interval
		err := callback()
		if err == nil {
			return nil
		}
		classification := classifier.Classify(err)
		if classification == pgerrors.NonRetriable {
			return err
		}
		if i == maxRetries {
			lastErr = err
		}
	}
	return lastErr
}
func (s *Storage) Close() error {

	return s.postgre.Close()
}

func (s *Storage) CheckConnection() error {
	const maxRetries = 3
	var lastErr error
	const interval = 2
	LinearBackoff := 1
	classifier := pgerrors.NewPostgresErrorClassifier()
	for i := 0; i <= maxRetries; i++ {
		time.Sleep(time.Duration(LinearBackoff) * time.Second)
		LinearBackoff += interval
		err := s.postgre.CheckConnection()
		if err == nil {
			return nil
		}
		classification := classifier.Classify(err)
		if classification == pgerrors.NonRetriable {
			return err
		}
		if i == maxRetries {
			lastErr = err
		}
	}
	return lastErr
}
func (s *Storage) SetGauge(key string, value float64, ctx context.Context) error {
	const maxRetries = 3
	var lastErr error
	const interval = 2
	LinearBackoff := 1
	classifier := pgerrors.NewPostgresErrorClassifier()
	for i := 0; i <= maxRetries; i++ {
		time.Sleep(time.Duration(LinearBackoff) * time.Second)
		LinearBackoff += interval
		err := s.postgre.SetGauge(key, value, ctx)
		if err == nil {
			return nil
		}
		classification := classifier.Classify(err)
		if classification == pgerrors.NonRetriable {
			return err
		}
		if i == maxRetries {
			lastErr = err
		}
	}
	return lastErr

}
func (s *Storage) SetCounter(key string, value int64, ctx context.Context) error {
	const maxRetries = 3
	var lastErr error
	const interval = 2
	LinearBackoff := 1
	classifier := pgerrors.NewPostgresErrorClassifier()
	for i := 0; i <= maxRetries; i++ {
		time.Sleep(time.Duration(LinearBackoff) * time.Second)
		LinearBackoff += interval
		err := s.postgre.SetCounter(key, value, ctx)
		if err == nil {
			return nil
		}
		classification := classifier.Classify(err)
		if classification == pgerrors.NonRetriable {
			return err
		}
		if i == maxRetries {
			lastErr = err
		}
	}
	return lastErr

}

func (s *Storage) GetAllGauges() (map[string]float64, error) {
	const maxRetries = 3
	var lastErr error
	const interval = 2
	LinearBackoff := 1
	classifier := pgerrors.NewPostgresErrorClassifier()
	for i := 0; i <= maxRetries; i++ {
		time.Sleep(time.Duration(LinearBackoff) * time.Second)
		LinearBackoff += interval
		m, err := s.postgre.GetAllGauges()
		if err == nil {
			return m, nil
		}
		classification := classifier.Classify(err)
		if classification == pgerrors.NonRetriable {
			return nil, err
		}
		if i == maxRetries {
			lastErr = err
		}
	}
	return nil, lastErr
}

func (s *Storage) GetAllCounters() (map[string]int64, error) {
	const maxRetries = 3
	var lastErr error
	const interval = 2
	LinearBackoff := 1
	classifier := pgerrors.NewPostgresErrorClassifier()
	for i := 0; i <= maxRetries; i++ {
		time.Sleep(time.Duration(LinearBackoff) * time.Second)
		LinearBackoff += interval
		m, err := s.postgre.GetAllCounters()
		if err == nil {
			return m, nil
		}
		classification := classifier.Classify(err)
		if classification == pgerrors.NonRetriable {
			return nil, err
		}
		if i == maxRetries {
			lastErr = err
		}
	}
	return nil, lastErr
}
