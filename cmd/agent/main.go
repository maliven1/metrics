package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/maliven1/metrics/internal/agent"
	"github.com/maliven1/metrics/internal/config"
	agenthandlers "github.com/maliven1/metrics/internal/handler/agent_handlers"
	"github.com/maliven1/metrics/internal/logger"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/storage"

	_ "net/http/pprof"
)

func main() {
	log, err := logger.Initialize()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	cfg := config.NewEnvAgentConfig()

	memStorage := storage.NewMemStorage()
	repo := repository.NewCache(memStorage, false, nil)
	service := agent.NewAgent(repo, cfg)
	client := agenthandlers.NewSendClient(service, cfg)

	var wg sync.WaitGroup
	wg.Add(2)
	go client.SendClientJSONMetrics(log, &wg)
	go client.SendClientBatchMetrics(log, &wg)
	go func() {
		http.ListenAndServe("localhost:6061", nil)
	}()
	wg.Wait()
}
