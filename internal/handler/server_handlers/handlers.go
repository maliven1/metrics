package serverhandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	models "github.com/maliven1/metrics/internal/model"
)

type Service interface {
	CheckAddPath(pathSplit []string) int
	GetMetric(pathSplit []string) (string, int)
	GetAllMetrics() (map[string]int64, map[string]float64)
	AddStructMetric(metric models.Metrics) int
	GetStructMetric(metric models.Metrics) (models.Metrics, int)
}

type AddHandler struct {
	AddHandler Service
}

func NewAddHandler(s Service) *AddHandler {
	return &AddHandler{AddHandler: s}
}

func (h AddHandler) GetBodyMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Add("content-type", "charset=utf-8")
		var buf bytes.Buffer
		var metric models.Metrics
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(models.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
			w.WriteHeader(models.StatusBadRequest)
			return
		}
		res, status := h.AddHandler.GetStructMetric(metric)
		resp, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(models.StatusInternalServerError)
			return
		}
		if status == models.StatusOK {
			w.WriteHeader(status)

			w.Write(resp)

			return
		}
		w.WriteHeader(status)
	}
}

func (h AddHandler) PostBodyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		var metric models.Metrics
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(models.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
			w.WriteHeader(models.StatusBadRequest)
			return
		}
		status := h.AddHandler.AddStructMetric(metric)
		resp, err := json.Marshal(metric)
		if err != nil {
			w.WriteHeader(models.StatusInternalServerError)
			return
		}
		w.Write(resp)
		w.WriteHeader(status)
		w.Header().Set("content-type", "text/plain")
		w.Header().Add("content-type", "charset=utf-8")
	}
}

func (h AddHandler) PostURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pathSplit := strings.Split(r.URL.Path, "/")
		status := h.AddHandler.CheckAddPath(pathSplit)
		w.WriteHeader(status)
		w.Header().Set("content-type", "text/plain")
		w.Header().Add("content-type", "charset=utf-8")
	}
}

func (h AddHandler) GetMetricHandler() http.HandlerFunc {
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
func (h AddHandler) GetAllMetricsHandler() http.HandlerFunc {
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
