package serverhandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	models "github.com/maliven1/metrics/internal/model"
	"go.uber.org/zap"
)

func (h Handler) PostMetricsHandler(log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var buf bytes.Buffer
		metrics := []models.Metrics{}
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if metrics == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.Handler.SetMetrics(metrics)

		w.WriteHeader(http.StatusOK)
	}
}

func (h Handler) PostBodyHandler(log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var buf bytes.Buffer
		var metric models.Metrics
		_, err := buf.ReadFrom(r.Body)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)

			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {

			w.WriteHeader(http.StatusBadRequest)

			return
		}
		err = h.Handler.AddStructMetric(metric)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp, err := json.Marshal(metric)
		if err != nil {
			log.Error("Marshal err: ", err, "status code: ", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	}
}

func (h Handler) PostURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		pathSplit := strings.Split(r.URL.Path, "/")
		if len(pathSplit) != 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err := h.Handler.CheckAddPath(pathSplit)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
