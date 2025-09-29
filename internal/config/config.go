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
)

type AgentConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

func ParseServerFlags() {
	flag.StringVar(&flagServerRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}
func ParseAgentFlags() {
	flag.StringVar(&flagAgentRunAddr, "a", ":8080", "address and port to run server")
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
