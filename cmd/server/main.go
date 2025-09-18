package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/service"
	"github.com/maliven1/metrics/internal/storage"
)

func main() {
	memStorage := storage.NewMemStorage()
	cache := repository.NewCache(memStorage)
	service := service.NewService(cache)
	h := serverhandlers.NewAddHandler(service)

	router := chi.NewRouter()
	router.Post(`/update/*`, h.PostHandler())
	router.Get(`/value/*`, h.GetMetricHandler())
	router.Get(`/`, h.GetAllMetricsHandler())
	log.Println("serv start")
	err := http.ListenAndServe(`:8080`, router)
	if err != nil {
		panic(err)
	}
}
