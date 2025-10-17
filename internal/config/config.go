package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

var (
	flagServerRunAddr string
	flagAgentRunAddr  string
	pollInterval      int
	reportInterval    int
	storeInterval     int
	fileStoragePath   string
	restore           bool
	postgreDSN        string
)

type AgentConfig struct {
	Address string `env:"ADDRESS"`
	//ReportInterval частота отправки метрик на сервер
	ReportInterval int `env:"REPORT_INTERVAL"`
	//PollInterval частота опроса метрик
	PollInterval int `env:"POLL_INTERVAL"`
}

type ServerConfig struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	//Restore нужно ли подгружать ранее сохраненные метрики в файле
	Restore    bool   `env:"RESTORE"`
	PostgreDSN string `env:"DATABASE_DSN"`
}

func parseServerFlags() {
	flag.StringVar(&flagServerRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&storeInterval, "i", 300, "the time interval after which the current server readings are saved")
	flag.StringVar(&fileStoragePath, "f", "", "path to the file where the current values are saved")
	flag.BoolVar(&restore, "r", false, "determines whether previously saved values from the specified file should be loaded when the server starts")
	flag.StringVar(&postgreDSN, "d", "postgres://postgres:12345678@localhost:5432/metrics_db?sslmode=disable", "postgres DSN")
	flag.Parse()
}
func parseAgentFlags() {
	flag.StringVar(&flagAgentRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.IntVar(&pollInterval, "p", 2, "metrics polling frequency")
	flag.Parse()

}

func NewEnvServerConfig() *ServerConfig {
	parseServerFlags()
	var cfg ServerConfig

	env.Parse(&cfg)
	if cfg.Address == "" {
		cfg.Address = flagServerRunAddr
	}
	if cfg.StoreInterval == 0 {
		cfg.StoreInterval = storeInterval
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = fileStoragePath
	}
	if !cfg.Restore {
		cfg.Restore = restore
	}
	if cfg.PostgreDSN == "" {
		cfg.PostgreDSN = postgreDSN
	}
	return &cfg

}

func NewEnvAgentConfig() *AgentConfig {
	parseAgentFlags()
	var cfg AgentConfig

	env.Parse(&cfg)
	if cfg.Address == "" {

		cfg.Address = flagAgentRunAddr
	}
	if cfg.PollInterval == 0 {

		cfg.PollInterval = pollInterval
	}
	if cfg.ReportInterval == 0 {

		cfg.ReportInterval = reportInterval
	}

	return &cfg
}
