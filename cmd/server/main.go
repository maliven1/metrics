package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/maliven1/metrics/internal/config"
	serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/service"
	"github.com/maliven1/metrics/internal/storage"
)

func main() {
	cfg := config.NewEnvServerConfig()
	memStorage := storage.NewMemStorage()
	cache := repository.NewCache(memStorage)
	service := service.NewService(cache)
	h := serverhandlers.NewAddHandler(service)

	router := chi.NewRouter()
	router.Post(`/update/*`, h.PostHandler())
	router.Get(`/value/*`, h.GetMetricHandler())
	router.Get(`/`, h.GetAllMetricsHandler())

	log.Println("serv start on", cfg.Address)
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

	log.Println("Shutdown Server ...")

	// Завершаем сервер с таймаутом
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	shutdownErr := srv.Shutdown(ctx)
	if shutdownErr != nil {
		log.Fatal("Server forced to shutdown:", shutdownErr)
	} else {
		select {
		case <-ctx.Done():
			log.Println("timeout of 30sec occurred")
		default:
			time.Sleep(time.Second * 5)
		}
	}

	log.Println("Server exited gracefully")
}
