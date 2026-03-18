// Package tests содержит примеры использования HTTP-хендлеров метрик.
// Демонстрирует основные сценарии работы с API системы мониторинга.
package serverhandlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"

	"log"

	"go.uber.org/zap"

	"github.com/go-chi/chi"
	serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers"
	models "github.com/maliven1/metrics/internal/model"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/service"
	"github.com/maliven1/metrics/internal/storage"

	"net/http"
)

// createTestServer создает тестовый HTTP сервер с настроенными маршрутами.
func createTestServer(handler *serverhandlers.Handler) *httptest.Server {
	r := chi.NewRouter()

	logger := zap.NewNop().Sugar()

	r.Get(`/`, handler.GetAllMetricsHandler())

	r.Post(`/value/`, handler.GetBodyMetricHandler(logger))
	r.Post(`/updates/`, handler.PostMetricsHandler(logger))
	r.Post(`/update/*`, handler.PostURLHandler())
	r.Post(`/update/`, handler.PostBodyHandler(logger))

	r.Get(`/value/*`, handler.GetMetricHandler())
	r.Get(`/ping`, handler.PingHandler(logger))

	return httptest.NewServer(r)
}

// Пример пакетного обновления метрик.
// Демонстрирует использование эндпоинта /updates для массового обновления.
func ExampleHandler_PostMetricsHandler() {
	//Инициализация зависимостей

	memStorage := storage.NewMemStorage()
	repo := repository.NewStorage(nil)
	cache := repository.NewCache(memStorage, false, nil)

	postgreService := service.NewPostgreService(repo, cache)
	logic := service.NewService(cache)
	h := serverhandlers.NewHandler(logic, postgreService)

	// Создание тестового сервера
	server := createTestServer(h)
	defer server.Close()

	// Отправка запроса на обновление метрик
	metrics := []models.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: func() *float64 { v := 12.6; return &v }(),
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Delta: func() *int64 { d := int64(1); return &d }(),
		},
	}

	jsonData, _ := json.Marshal(metrics)

	resp, err := http.Post(
		server.URL+"/updates/",
		"application/json",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		log.Fatalf("Ошибка отправки запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check status code before processing response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 200

}

// Пример получения всех метрик.
// Демонстрирует использование эндпоинта / для получения списка всех метрик.
func ExampleHandler_GetAllMetricsHandler() {
	// Инициализация зависимостей
	memStorage := storage.NewMemStorage()
	repo := repository.NewStorage(nil)
	cache := repository.NewCache(memStorage, false, nil)

	postgreService := service.NewPostgreService(repo, cache)
	logic := service.NewService(cache)
	h := serverhandlers.NewHandler(logic, postgreService)

	// Создание тестового сервера
	server := createTestServer(h)
	defer server.Close()

	// Отправка запроса на получение всех метрик
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		log.Fatalf("Ошибка отправки запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check status code before processing response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 200
}

// Пример добавления метрики через JSON тело.
// Демонстрирует использование эндпоинта /update/ для добавления метрики.
func ExampleHandler_PostBodyHandler() {
	// Инициализация зависимостей
	memStorage := storage.NewMemStorage()
	repo := repository.NewStorage(nil)
	cache := repository.NewCache(memStorage, false, nil)

	postgreService := service.NewPostgreService(repo, cache)
	logic := service.NewService(cache)
	h := serverhandlers.NewHandler(logic, postgreService)

	// Создание тестового сервера
	server := createTestServer(h)
	defer server.Close()

	// Подготовка метрики
	metric := models.Metrics{
		ID:    "TestGauge",
		MType: "gauge",
		Value: func() *float64 { v := 42.0; return &v }(),
	}

	jsonData, _ := json.Marshal(metric)

	resp, err := http.Post(
		server.URL+"/update/",
		"application/json",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		log.Fatalf("Ошибка отправки запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check status code before processing response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 200
}

// Пример обновления метрики через URL параметры.
// Демонстрирует использование эндпоинта /update/{type}/{name}/{value}.
func ExampleHandler_PostURLHandler() {
	// Инициализация зависимостей
	memStorage := storage.NewMemStorage()
	repo := repository.NewStorage(nil)
	cache := repository.NewCache(memStorage, false, nil)

	postgreService := service.NewPostgreService(repo, cache)
	logic := service.NewService(cache)
	h := serverhandlers.NewHandler(logic, postgreService)

	// Создание тестового сервера
	server := createTestServer(h)
	defer server.Close()

	// Отправка запроса на обновление метрики gauge
	resp, err := http.Post(
		server.URL+"/update/gauge/TestGauge/123.45",
		"text/plain",
		nil,
	)
	if err != nil {
		log.Fatalf("Ошибка отправки запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check status code before processing response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 200
}

// Пример получения метрики через JSON тело.
// Демонстрирует использование эндпоинта /value/ для получения метрики.
func ExampleHandler_GetBodyMetricHandler() {
	// Инициализация зависимостей
	memStorage := storage.NewMemStorage()
	repo := repository.NewStorage(nil)
	cache := repository.NewCache(memStorage, false, nil)

	postgreService := service.NewPostgreService(repo, cache)
	logic := service.NewService(cache)
	h := serverhandlers.NewHandler(logic, postgreService)

	// Создание тестового сервера
	server := createTestServer(h)
	defer server.Close()

	// Сначала добавляем метрику, чтобы она существовала
	metric := models.Metrics{
		ID:    "TestGauge",
		MType: "gauge",
		Value: func() *float64 { v := 99.9; return &v }(),
	}
	jsonData, _ := json.Marshal(metric)
	res, err := http.Post(
		server.URL+"/update/",
		"application/json",
		bytes.NewReader(jsonData),
	)

	if err != nil {
		log.Fatalf("Ошибка добавления метрики: %v\n", err)
		return
	}
	defer res.Body.Close()
	// Теперь запрашиваем её
	reqMetric := models.Metrics{
		ID:    "TestGauge",
		MType: "gauge",
	}
	reqData, _ := json.Marshal(reqMetric)

	resp, err := http.Post(
		server.URL+"/value/",
		"application/json",
		bytes.NewReader(reqData),
	)
	if err != nil {
		log.Fatalf("Ошибка отправки запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check status code before processing response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 200
}

// Пример получения метрики через URL параметры.
// Демонстрирует использование эндпоинта /value/{type}/{id}.
func ExampleHandler_GetMetricHandler() {
	// Инициализация зависимостей
	memStorage := storage.NewMemStorage()
	repo := repository.NewStorage(nil)
	cache := repository.NewCache(memStorage, false, nil)

	postgreService := service.NewPostgreService(repo, cache)
	logic := service.NewService(cache)
	h := serverhandlers.NewHandler(logic, postgreService)

	// Создание тестового сервера
	server := createTestServer(h)
	defer server.Close()

	// Сначала добавляем метрику, чтобы она существовала
	metric := models.Metrics{
		ID:    "TestCounter",
		MType: "counter",
		Delta: func() *int64 { d := int64(5); return &d }(),
	}
	jsonData, _ := json.Marshal(metric)
	res, err := http.Post(
		server.URL+"/update/",
		"application/json",
		bytes.NewReader(jsonData),
	)

	if err != nil {
		log.Fatalf("Ошибка добавления метрики: %v\n", err)
		return
	}
	defer res.Body.Close()
	// Теперь запрашиваем её через URL
	resp, err := http.Get(server.URL + "/value/counter/TestCounter")
	if err != nil {
		log.Fatalf("Ошибка отправки запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check status code before processing response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 200
}

// Пример проверки соединения с базой данных.
// Демонстрирует использование эндпоинта /ping.
func ExampleHandler_PingHandler() {
	// Инициализация зависимостей
	memStorage := storage.NewMemStorage()
	repo := repository.NewStorage(nil)
	cache := repository.NewCache(memStorage, false, nil)

	postgreService := service.NewPostgreService(repo, cache)
	logic := service.NewService(cache)
	h := serverhandlers.NewHandler(logic, postgreService)

	// Создание тестового сервера
	server := createTestServer(h)
	defer server.Close()

	// Отправка запроса на проверку соединения
	resp, err := http.Get(server.URL + "/ping")
	if err != nil {
		log.Fatalf("Ошибка отправки запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check status code before processing response
	if resp.StatusCode != http.StatusInternalServerError {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 500
}
