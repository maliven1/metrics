// Package main
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

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func defaultIfEmpty(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}

func main() {
	fmt.Printf("Build version: %s\n", defaultIfEmpty(buildVersion))
	fmt.Printf("Build date: %s\n", defaultIfEmpty(buildDate))
	fmt.Printf("Build commit: %s\n", defaultIfEmpty(buildCommit))
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
	wg.Add(3)
	go client.SendClientJSONMetrics(log, &wg)
	go client.SendClientBatchMetrics(log, &wg)
	go func() {
		http.ListenAndServe("localhost:6061", nil)
	}()
	wg.Wait()
}
