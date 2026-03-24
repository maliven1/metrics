// Package serverhandlers
package serverhandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	models "github.com/maliven1/metrics/internal/model"
	"go.uber.org/zap"
)

// PostMetricsHandler godoc
// @Tags Info
// @Summary Пакетное обновление метрик
// @Description Принимает массив метрик в JSON формате и обновляет их все за один запрос
// @Accept json
// @Produce json
// @Param metrics body []models.Metrics true "Массив метрик для обновления"
// @Success 200 {object} map[string]string "Пример: {\"status\":\"OK\"}"
// @Failure 400 {object} map[string]string "Пример: {\"status\":\"StatusBadRequest\"}"
// @Router /updates/ [post]
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
			return
		}
		if metrics == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.Handler.SetMetrics(metrics)

		w.WriteHeader(http.StatusOK)
	}
}

// PostBodyHandler godoc
// @Tags Info
// @Summary Добавление метрики
// @Description Принимает метрику в JSON формате и добавляет ее в базу данных
// @Accept json
// @Produce json
// @Param metrics body models.Metrics true "Метрика для добавления"
// @Success 200 {object} models.Metrics "Пример: {\"id\":\"1\",\"type\":\"gauge\",\"value\":1}"
// @Failure 400 {object} map[string]string "Пример: {\"status\":\"StatusBadRequest\"}"
// @Router /update/ [post]
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

// UpdateMetric godoc
// @Tags Info
// @Summary Обновление метрики
// @Description Обновляет метрику. Поддерживает два формата:
//  1. JSON (новый) - отправка объекта Metrics в теле запроса
//  2. URL параметры (старый) - /update/{type}/{name}/{value}
//
// @Accept json
// @Accept plain
// @Produce json
// @Produce plain
// @Param metric body models.Metrics false "Метрика в JSON формате"
// @Param type path string true "Тип метрики" Enums(gauge, counter)
// @Param name path string true "Имя метрики"
// @Param value path string true "Значение метрики"
// @Success 200 {object} map[string]string "Объект со статусом OK"
// @Success 200 {string} string "Строка со статусом OK"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 400 {string} string "Неверный запрос: bad gauge value, bad counter value, unknown metric type, metric ID is required, gauge value is required, counter delta is required"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Внутренняя ошибка сервера: store error"
// @Router /update/{type}/{name}/{value} [post]
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
