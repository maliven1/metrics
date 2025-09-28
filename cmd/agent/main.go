package main

import (
	"github.com/maliven1/metrics/internal/agent"
	"github.com/maliven1/metrics/internal/config"
	agenthandlers "github.com/maliven1/metrics/internal/handler/agent_handlers"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/storage"
)

func main() {
	cfg := config.NewEnvAgentConfig()
	memStorage := storage.NewMemStorage()
	cache := repository.NewCache(memStorage)
	service := agent.NewAgent(cache, cfg)
	client := agenthandlers.NewSendClient(service, cfg)
	go client.SendClientMetrics()
	client.SendClientJSONMetrics()

}
