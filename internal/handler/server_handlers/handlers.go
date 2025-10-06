package serverhandlers

import (
	models "github.com/maliven1/metrics/internal/model"
)

type Service interface {
	CheckAddPath(pathSplit []string) int
	GetMetric(pathSplit []string) (string, int)
	GetAllMetrics() (map[string]int64, map[string]float64)
	AddStructMetric(metric models.Metrics) int
	GetStructMetric(metric models.Metrics) (models.Metrics, int)
}
type PostgreService interface {
	CheckConnection() int
}

type Handler struct {
	Handler        Service
	PostgreHandler PostgreService
}

func NewHandler(s Service, p PostgreService) *Handler {
	return &Handler{Handler: s, PostgreHandler: p}
}
