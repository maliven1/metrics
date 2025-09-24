package agenthandlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/maliven1/metrics/internal/config"
	models "github.com/maliven1/metrics/internal/model"
)

type Agent interface {
	GetMetrics() (map[string]float64, map[string]int64)
	CollectMetrics()
}

type SendClient struct {
	AddHandler Agent
	cfg        *config.Config
}

func NewSendClient(s Agent, cfg *config.Config) *SendClient {
	return &SendClient{AddHandler: s, cfg: cfg}
}

func (s SendClient) SendClientMetrics() {
	endpoint := "http://" + s.cfg.Address + "/update/"
	client := &http.Client{}
	go s.AddHandler.CollectMetrics()
	for {
		time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
		Gauge, Counter := s.AddHandler.GetMetrics()
		for i, v := range Gauge {
			if i == "" {
				continue
			}

			data, _ := url.JoinPath(models.Gauge, i, fmt.Sprint(v))

			request, err := http.NewRequest(http.MethodPost, endpoint+data, nil)
			if err != nil {
				log.Println(err)
			}

			request.Header.Add("Content-Type", "Content-Type: text/plain")

			response, err := client.Do(request)
			if err != nil {
				log.Println(err)
			}

			response.Body.Close()
		}
		for i, v := range Counter {
			if i == "" {
				continue
			}

			data, _ := url.JoinPath(models.Counter, i, fmt.Sprint(v))

			request, err := http.NewRequest(http.MethodPost, endpoint+data, nil)
			if err != nil {
				log.Println(err)
			}

			request.Header.Add("Content-Type", "Content-Type: text/plain")

			response, err := client.Do(request)
			if err != nil {
				log.Println(err)
			}
			response.Body.Close()
		}
	}
}
