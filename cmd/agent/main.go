package main

import (
	"github.com/maliven1/metrics/internal/agent"
	agentHandlers "github.com/maliven1/metrics/internal/handler/agent_handlers"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/storage"
)

func main() {

	memStorage := storage.NewMemStorage()
	cache := repository.NewCache(memStorage)
	service := agent.NewAgent(cache)
	client := agentHandlers.NewSendClient(service)
	client.SendClientMetrics()

}
