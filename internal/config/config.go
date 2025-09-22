package config

import (
	"flag"

	models "github.com/maliven1/metrics/internal/model"
)

var (
	flagServerRunAddr string
	flagAgentRunAddr  string
)

func ParseServerFlags() string {
	flag.StringVar(&flagServerRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
	return flagServerRunAddr
}
func ParseAgentFlags() string {
	flag.StringVar(&flagAgentRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&models.ReportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.IntVar(&models.PollInterval, "p", 2, "metrics polling frequency")
	flag.Parse()
	return flagAgentRunAddr
}
