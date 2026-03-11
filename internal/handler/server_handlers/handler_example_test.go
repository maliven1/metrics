package serverhandlers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/go-chi/chi"
	"github.com/maliven1/metrics/internal/config"
	serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers"
	"github.com/maliven1/metrics/internal/logger"
	models "github.com/maliven1/metrics/internal/model"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/service"
	"github.com/maliven1/metrics/internal/storage"
	"go.uber.org/zap"

	"net/http"
)

func createTestServer(handler *serverhandlers.Handler, log *zap.SugaredLogger, cfg config.ServerConfig) *httptest.Server {
	r := chi.NewRouter()

	r.Get(`/`, handler.GetAllMetricsHandler())

	r.Post(`/value/`, handler.GetBodyMetricHandler(log))

	r.Post(`/updates/`, handler.PostMetricsHandler(log))
	r.Post(`/update/*`, handler.PostURLHandler())
	r.Post(`/update/`, handler.PostBodyHandler(log))

	r.Get(`/value/*`, handler.GetMetricHandler())
	r.Get(`/ping`, handler.PingHandler(log))

	return httptest.NewServer(r)
}

func ExampleHandler_PostMetricsHandler() {
	cfg := config.NewEnvServerConfig()
	log, err := logger.Initialize()
	if err != nil {
		log.Info(err)
		return
	}
	defer log.Sync()

	var usePostgreSQL bool

	postgreStorage, err := storage.NewPostgreDB(*cfg, log)
	if err != nil {
		log.Info(err)

	} else {
		usePostgreSQL = true
	}
	memStorage := storage.NewMemStorage()
	repo := repository.NewStorage(postgreStorage)
	cache := repository.NewCache(memStorage, usePostgreSQL, postgreStorage)

	postgreService := service.NewPostgreService(repo, cache)
	logic := service.NewService(cache)
	h := serverhandlers.NewHandler(logic, postgreService)
	server := createTestServer(h, log, *cfg)

	metrics := []models.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: func() *float64 { v := 1234.56; return &v }(),
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
		log.Errorf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Errorf("Failed to decode response: %v", err)
		return
	}

	json.NewDecoder(resp.Body).Decode(&response)
	log.Infof("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 200

}
