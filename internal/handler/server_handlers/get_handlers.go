package serverhandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	models "github.com/maliven1/metrics/internal/model"
	"go.uber.org/zap"
)

// GetBodyMetricHandler godoc
// @Tags Info
// @Summary Получение метрики по JSON
// @Description Принимает JSON с полями id и type, возвращает структуру метрики с полями ID, Type, Value.
// @Accept json
// @Produce json
// @Param metric body models.Metrics true "Метрика для поиска (нужны только id и type)"
// @Success 200 {object} models.Metrics "Метрика"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /value/ [get]
func (h Handler) GetBodyMetricHandler(log *zap.SugaredLogger) http.HandlerFunc {
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
		if metric.ID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := h.Handler.GetStructMetric(metric)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		resp, err := json.Marshal(res)
		if err != nil {
			log.Error("Marshal err: ", err, "status code: ", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

// GetMetricHandler godoc
// @Tags Info
// @Summary Получение метрики по URL
// @Description Принимает URL с полями id и type, возвращает структуру метрики с полями ID, Type, Value.
// @Accept plain
// @Produce plain
// @Param id path string true "ID метрики"
// @Param type path string true "Тип метрики"
// @Success 200 {object} models.Metrics "Метрика"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /value/ID/type [get]
func (h Handler) GetMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		pathSplit := strings.Split(r.URL.Path, "/")
		if len(pathSplit) != 4 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics, err := h.Handler.GetMetric(pathSplit)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.PlainText(w, r, metrics)
	}
}

// GetAllMetricsHandler godoc
// @Tags Info
// @Summary Получение всех метрик
// @Description Возвращает все метрики в формате HTML.
// @Produce html
// @Success 200 {string} string "Все метрики"
// @Failure 500 {string} string "Internal Server Error"
// @Router / [get]
func (h Handler) GetAllMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Header().Add("content-type", "charset=utf-8")
		count, gauge := h.Handler.GetAllMetrics()
		jsonCount, err := json.Marshal(count)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonGauge, err := json.Marshal(gauge)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		str := string(jsonCount) + string(jsonGauge)
		w.WriteHeader(http.StatusOK)

		render.HTML(w, r, str)

	}
}

// PingHandler godoc
// @Tags Info
// @Summary Проверка соединения с базой данных
// @Description Проверяет соединение с базой данных.
// @Produce plain
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /ping [get]
func (h Handler) PingHandler(log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain")
		err := h.PostgreHandler.CheckConnection()
		if err != nil {
			log.Info("status cod: ", http.StatusInternalServerError, "ping postgreDB failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
