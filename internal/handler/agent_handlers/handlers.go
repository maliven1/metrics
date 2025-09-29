package agenthandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/maliven1/metrics/internal/config"
	models "github.com/maliven1/metrics/internal/model"
	"go.uber.org/zap"
)

type Agent interface {
	GetMetrics() (map[string]float64, map[string]int64)
	CollectMetrics()
}

type SendClient struct {
	AddHandler Agent
	cfg        *config.AgentConfig
}

func NewSendClient(s Agent, cfg *config.AgentConfig) *SendClient {
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

func (s SendClient) SendClientJSONMetrics(log *zap.SugaredLogger) {
	time.Sleep(1 * time.Second)
	endpoint := "http://" + s.cfg.Address + "/update/"
	log.Info("start agent on endpoint: ", endpoint)
	client := &http.Client{}
	go s.AddHandler.CollectMetrics()
	for {
		time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
		Gauge, Counter := s.AddHandler.GetMetrics()
		for i, v := range Gauge {
			if i == "" {
				continue
			}

			metric := models.Metrics{MType: models.Gauge, ID: i, Value: &v}
			data, err := json.Marshal(metric)
			if err != nil {
				log.Info(err)
			}
			reader := bytes.NewReader(data)

			request, err := http.NewRequest(http.MethodPost, endpoint, reader)
			if err != nil {
				log.Info(err)
			}

			request.Header.Set("content-type", "application/json")

			response, err := client.Do(request)
			if err != nil {
				log.Info(err)
				continue
			}

			response.Body.Close()
		}
		for i, v := range Counter {
			if i == "" {
				continue
			}

			metric := models.Metrics{MType: models.Counter, ID: i, Delta: &v}
			data, err := json.Marshal(metric)
			if err != nil {
				log.Info(err)
			}
			reader := bytes.NewReader(data)
			request, err := http.NewRequest(http.MethodPost, endpoint, reader)
			if err != nil {
				log.Info(err)
			}

			request.Header.Set("content-type", "application/json")

			response, err := client.Do(request)
			if err != nil {
				log.Info(err)
				continue
			}
			response.Body.Close()
		}
	}
}
