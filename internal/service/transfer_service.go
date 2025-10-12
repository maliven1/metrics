package service

import (
	"time"

	"github.com/maliven1/metrics/internal/config"
	"github.com/maliven1/metrics/internal/repository"
)

func TransferCacheToPostgreSQL(memRepo MemRepo, postgreRepo repository.Postgre, cfg config.ServerConfig) {
	for {
		time.Sleep(time.Duration(cfg.StoreInterval) * time.Second)
		gauges := memRepo.GetGauge()
		for key, value := range gauges {
			postgreRepo.SetGauge(key, value)
		}

		counters := memRepo.GetCounter()
		for key, value := range counters {
			postgreRepo.SetCounter(key, value)
		}
	}

}
