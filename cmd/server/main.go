package main

import (
	"log"
	"net/http"

	"github.com/maliven1/metrics/internal/config"
	"github.com/maliven1/metrics/internal/handler"
)

func main() {
	memStorage := config.InitMemStorage()
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, handler.PostHandler(memStorage))
	log.Println("serv start")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
