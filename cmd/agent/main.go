package main

import (
	"fmt"

	"github.com/maliven1/metrics/internal/agent"
	"github.com/maliven1/metrics/internal/config"
	agenthandlers "github.com/maliven1/metrics/internal/handler/agent_handlers"
	"github.com/maliven1/metrics/internal/logger"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/storage"
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
	cache := repository.NewCache(memStorage)
	service := agent.NewAgent(cache, cfg)
	client := agenthandlers.NewSendClient(service, cfg)

	client.SendClientJSONMetrics(log)

}
