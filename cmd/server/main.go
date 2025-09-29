package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/maliven1/metrics/internal/config"
	serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers"
	"github.com/maliven1/metrics/internal/logger"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/service"
	"github.com/maliven1/metrics/internal/storage"
)

func main() {
	cfg := config.NewEnvServerConfig()
	log := logger.Initialize()
	memStorage := storage.NewMemStorage()
	cache := repository.NewCache(memStorage)
	service := service.NewService(cache)
	h := serverhandlers.NewAddHandler(service)

	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Use(logger.WithLogging)
		r.Post(`/update/*`, h.PostURLHandler())
		r.Post(`/value/`, h.GetBodyMetricHandler(log))
		r.Post(`/update/`, h.PostBodyHandler())
		r.Get(`/value/*`, h.GetMetricHandler())
		r.Get(`/`, h.GetAllMetricsHandler())

	})
	log.Info("serv start on ", cfg.Address, " time:", time.Now())
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutdown Server ...", " time:", time.Now())

	// Завершаем сервер с таймаутом
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	shutdownErr := srv.Shutdown(ctx)
	if shutdownErr != nil {
		log.Fatal("Server forced to shutdown:", shutdownErr)
	} else {
		select {
		case <-ctx.Done():
			log.Info("timeout of 30sec occurred", " time:", time.Now())
		default:
			time.Sleep(time.Second * 5)
		}
	}

	log.Info("Server exited gracefully", " time:", time.Now())
}
