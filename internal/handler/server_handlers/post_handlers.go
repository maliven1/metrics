package serverhandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	models "github.com/maliven1/metrics/internal/model"
	"go.uber.org/zap"
)

func (h Handler) PostMetricsHandler(ctx context.Context, log *zap.SugaredLogger) http.HandlerFunc {
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
		status, err := h.PostgreHandler.SetMetrics(metrics, ctx)
		if err != nil {
			log.Error(err)
		}
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

			w.WriteHeader(http.StatusBadRequest)

			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {

			w.WriteHeader(http.StatusBadRequest)

			return
		}
		status := h.Handler.AddStructMetric(metric)
		resp, err := json.Marshal(metric)
		if err != nil {
			log.Error("Marshal err: ", err, "status code: ", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
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
