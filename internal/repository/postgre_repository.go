// Package repository
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
)

type Storage struct {
	postgre Postgre
}

func NewStorage(postgre Postgre) *Storage {
	return &Storage{postgre: postgre}
}

func (s *Storage) Close() error {

	return s.postgre.Close()
}

func (s *Storage) CheckConnection() error {
	if !usePostgre {
		return fmt.Errorf("postgres is not used")
	}
	var delay = time.Second           // Начальная задержка
	const increment = 2 * time.Second // Увеличение задержки на 2 секунды после каждой попытки

	err := retry.Do(func() error {
		return s.postgre.CheckConnection()
	}, retry.Attempts(3), retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
		if n > 0 {
			delay += increment
		}
		return delay
	}))
	return err
}
func (s *Storage) SetGauge(key string, value float64, ctx context.Context) error {
	var delay = time.Second           // Начальная задержка
	const increment = 2 * time.Second // Увеличение задержки на 2 секунды после каждой попытки

	err := retry.Do(func() error {
		return s.postgre.SetGauge(key, value, ctx)
	}, retry.Attempts(3), retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
		if n > 0 {
			delay += increment
		}
		return delay
	}), retry.Context(ctx))
	return err

}
func (s *Storage) SetCounter(key string, value int64, ctx context.Context) error {
	var delay = time.Second           // Начальная задержка
	const increment = 2 * time.Second // Увеличение задержки на 2 секунды после каждой попытки

	err := retry.Do(func() error {
		return s.postgre.SetCounter(key, value, ctx)
	}, retry.Attempts(3), retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
		if n > 0 {
			delay += increment
		}
		return delay
	}), retry.Context(ctx))
	return err

}

func (s *Storage) GetAllGauges() (map[string]float64, error) {
	var delay = time.Second           // Начальная задержка
	const increment = 2 * time.Second // Увеличение задержки на 2 секунды после каждой попытки
	m := make(map[string]float64)
	err := retry.Do(func() error {
		var err error
		m, err = s.postgre.GetAllGauges()
		return err
	}, retry.Attempts(3), retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
		if n > 0 {
			delay += increment
		}
		return delay
	}))
	return m, err

}

func (s *Storage) GetAllCounters() (map[string]int64, error) {
	var delay = time.Second           // Начальная задержка
	const increment = 2 * time.Second // Увеличение задержки на 2 секунды после каждой попытки
	m := make(map[string]int64)
	err := retry.Do(func() error {
		var err error
		m, err = s.postgre.GetAllCounters()
		return err
	}, retry.Attempts(3), retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
		if n > 0 {
			delay += increment
		}
		return delay
	}))
	return m, err
}
