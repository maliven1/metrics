package main

import (
	"github.com/maliven1/metrics/internal/app"
)

// @title Metrics API
// @version 1.0
// @description API для работы с метриками системы
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	app.Run()
}
