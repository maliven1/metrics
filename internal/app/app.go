// Package app
package app

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
	"github.com/maliven1/metrics/internal/router"
	"github.com/maliven1/metrics/internal/service"
	"github.com/maliven1/metrics/internal/storage"

	_ "net/http/pprof"
)

func Run() {

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
	cahce := repository.NewCache(memStorage, usePostgreSQL, postgreStorage)

	postgreService := service.NewPostgreService(repo, cahce)
	logic := service.NewService(cahce)
	h := serverhandlers.NewHandler(logic, postgreService)

	go logic.InitFile(*cfg, log)

	r := chi.NewRouter()
	router.NewRouter(r, h, log, *cfg)

	log.Info("serv start on ", cfg.Address, " time:", time.Now())
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// listen on localhost:6060 for pprof
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutdown Server ...", " time:", time.Now())
	// Завершаем сервер с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	shutdownErr := srv.Shutdown(ctx)
	if shutdownErr != nil {
		log.Fatal("Server forced to shutdown:", shutdownErr)
	}

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			log.Info("timeout of 30sec occurred", " time:", time.Now())
		}
	default:

	}

	log.Info("Server exited gracefully", " time:", time.Now())
}
