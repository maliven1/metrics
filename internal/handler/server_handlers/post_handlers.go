package serverhandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	models "github.com/maliven1/metrics/internal/model"
	"go.uber.org/zap"
)

func (h Handler) PostMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var buf bytes.Buffer
		metrics := []models.Metrics{}
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(models.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
			w.WriteHeader(models.StatusBadRequest)
		}
		status := h.PostgreHandler.SetMetrics(metrics)
		w.WriteHeader(status)
	}
}

func (h Handler) PostBodyHandler(log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
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
		status := h.Handler.AddStructMetric(metric)
		resp, err := json.Marshal(metric)
		if err != nil {
			log.Error("Marshal err: ", err, "status code: ", models.StatusInternalServerError)
			w.WriteHeader(models.StatusInternalServerError)
			return
		}
		w.WriteHeader(status)
		w.Write(resp)

	}
}

func (h Handler) PostURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pathSplit := strings.Split(r.URL.Path, "/")
		status := h.Handler.CheckAddPath(pathSplit)
		w.WriteHeader(status)
		w.Header().Set("content-type", "text/plain")
		w.Header().Add("content-type", "charset=utf-8")
	}
}
