package config

import (
	"flag"

	models "github.com/maliven1/metrics/internal/model"
)

func ParseServerFlags() {
	flag.StringVar(&models.FlagServerRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
}

func ParseAgentFlags() {
	flag.StringVar(&models.FlagAgentRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&models.ReportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.IntVar(&models.PollInterval, "p", 2, "metrics polling frequency")
	flag.Parse()
}
