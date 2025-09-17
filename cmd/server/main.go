package main

import (
	"log"
	"net/http"

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

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, h.PostHandler())
	log.Println("serv start")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
