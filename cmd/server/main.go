// Package main
package main

import (
	"fmt"

	"github.com/maliven1/metrics/internal/app"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func defaultIfEmpty(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}

// @title Metrics API
// @version 1.0
// @description API для работы с метриками системы
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	fmt.Printf("Build version: %s\n", defaultIfEmpty(buildVersion))
	fmt.Printf("Build date: %s\n", defaultIfEmpty(buildDate))
	fmt.Printf("Build commit: %s\n", defaultIfEmpty(buildCommit))

	app.Run()
}
