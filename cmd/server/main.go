package main

import (
	"log"
	"net/http"

	"github.com/maliven1/metrics/internal/config"
	serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers"
	models "github.com/maliven1/metrics/internal/model"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/router"
	"github.com/maliven1/metrics/internal/service"
	"github.com/maliven1/metrics/internal/storage"
)

func main() {
	config.ParseServerFlags()
	memStorage := storage.NewMemStorage()
	cache := repository.NewCache(memStorage)
	service := service.NewService(cache)
	h := serverhandlers.NewAddHandler(service)
	r := router.New(*h)

	log.Println("serv start on", models.FlagServerRunAddr)
	err := http.ListenAndServe(models.FlagServerRunAddr, r)
	if err != nil {
		panic(err)
	}
}
