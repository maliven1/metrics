package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maliven1/metrics/internal/config"
	serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers"
	models "github.com/maliven1/metrics/internal/model"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/service"
	"github.com/maliven1/metrics/internal/storage"
)

func main() {
	config.ParseServerFlags()
	memStorage := storage.NewMemStorage()
	cache := repository.NewCache(memStorage)
	service := service.NewService(cache)
	h := serverhandlers.NewAddHandler(service)

	router := chi.NewRouter()
	router.Post(`/update/*`, h.PostHandler())
	router.Get(`/value/*`, h.GetMetricHandler())
	router.Get(`/`, h.GetAllMetricsHandler())
	log.Println("serv start on", models.FlagRunAddr)
	err := http.ListenAndServe(models.FlagRunAddr, router)
	if err != nil {
		panic(err)
	}
}
