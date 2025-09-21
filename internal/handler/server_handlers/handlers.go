package serverhandlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	models "github.com/maliven1/metrics/internal/model"
)

//go:generate mockgen -source=./handlers.go -destination=mocks/mock.go
type Service interface {
	CheckAddPath(pathSplit []string) int
	GetMetric(pathSplit []string) (string, int)
	GetAllMetrics() (map[string]int64, map[string]float64)
}

type AddHandler struct {
	AddHandler Service
}

func NewAddHandler(s Service) *AddHandler {
	return &AddHandler{AddHandler: s}
}

func PostHandler(h AddHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pathSplit := strings.Split(r.URL.Path, "/")
		status := h.AddHandler.CheckAddPath(pathSplit)
		w.WriteHeader(status)
		w.Header().Set("content-type", "text/plain")
		w.Header().Add("content-type", "charset=utf-8")
	}
}

func GetMetricHandler(h AddHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain")
		w.Header().Add("content-type", "charset=utf-8")
		pathSplit := strings.Split(r.URL.Path, "/")
		metrics, status := h.AddHandler.GetMetric(pathSplit)
		if status == models.StatusOK {
			w.WriteHeader(status)

			render.PlainText(w, r, metrics)

			return
		}
		w.WriteHeader(status)

	}
}

func GetAllMetricsHandler(h AddHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Header().Add("content-type", "charset=utf-8")
		count, gauge := h.AddHandler.GetAllMetrics()
		jsonCount, err := json.Marshal(count)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		jsonGauge, err := json.Marshal(gauge)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		str := string(jsonCount) + string(jsonGauge)
		w.WriteHeader(models.StatusOK)

		render.HTML(w, r, str)

	}
}
