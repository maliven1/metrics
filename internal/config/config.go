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
)

type AgentConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

type ServerConfig struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func ParseServerFlags() {
	flag.StringVar(&flagServerRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&storeInterval, "i", 300, "the time interval after which the current server readings are saved")
	flag.StringVar(&fileStoragePath, "f", "history", "path to the file where the current values are saved")
	flag.BoolVar(&restore, "r", false, "determines whether previously saved values from the specified file should be loaded when the server starts")
	flag.Parse()
}
func ParseAgentFlags() {
	flag.StringVar(&flagAgentRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.IntVar(&pollInterval, "p", 2, "metrics polling frequency")
	flag.Parse()

}

func NewEnvServerConfig() *ServerConfig {
	ParseServerFlags()
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
	return &cfg

}

func NewEnvAgentConfig() *AgentConfig {
	ParseAgentFlags()
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
