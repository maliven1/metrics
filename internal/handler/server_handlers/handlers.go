package serverhandlers

import (
	models "github.com/maliven1/metrics/internal/model"
)

type Service interface {
	CheckAddPath(pathSplit []string) error
	GetMetric(pathSplit []string) (string, error)
	GetAllMetrics() (map[string]int64, map[string]float64)
	AddStructMetric(metric models.Metrics) error
	GetStructMetric(metric models.Metrics) (models.Metrics, error)
	SetMetrics(metrics []models.Metrics)
}
type PostgreService interface {
	CheckConnection() error
}

type Handler struct {
	Handler        Service
	PostgreHandler PostgreService
}

func NewHandler(s Service, p PostgreService) *Handler {
	return &Handler{Handler: s, PostgreHandler: p}
}
