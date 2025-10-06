package agenthandlers

import (
	"bytes"
	"compress/gzip"
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
		gauge, counter := s.AddHandler.GetMetrics()
		for i, v := range gauge {
			if i == "" {
				return
			}

			data, _ := url.JoinPath(models.Gauge, i, fmt.Sprint(v))

			request, err := http.NewRequest(http.MethodPost, endpoint+data, nil)
			if err != nil {
				log.Println(err)
				return
			}

			request.Header.Set("Content-Type", "Content-Type: text/plain")

			response, err := client.Do(request)
			if err != nil {
				log.Println(err)
				continue
			}

			defer response.Body.Close()
		}
		for i, v := range counter {
			if i == "" {
				return
			}

			data, _ := url.JoinPath(models.Counter, i, fmt.Sprint(v))

			request, err := http.NewRequest(http.MethodPost, endpoint+data, nil)
			if err != nil {
				log.Println(err)
				return
			}

			request.Header.Set("Content-Type", "Content-Type: text/plain")

			response, err := client.Do(request)
			if err != nil {
				log.Println(err)
				continue
			}
			defer response.Body.Close()
		}
	}
}

func (s SendClient) SendClientJSONMetrics(log *zap.SugaredLogger) {

	endpoint := "http://" + s.cfg.Address + "/update/"
	log.Info("start agent on endpoint: ", endpoint)
	client := &http.Client{}
	go s.AddHandler.CollectMetrics()
	for {
		time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
		gauge, counter := s.AddHandler.GetMetrics()
		for i, v := range gauge {
			if i == "" {
				return
			}

			metric := models.Metrics{MType: models.Gauge, ID: i, Value: &v}
			data, err := json.Marshal(metric)
			if err != nil {
				log.Error(err)
				return
			}

			var buf bytes.Buffer
			gzipWriter := gzip.NewWriter(&buf)

			_, err = gzipWriter.Write(data)
			if err != nil {
				log.Info(err)
				return
			}
			err = gzipWriter.Flush()
			if err != nil {
				log.Info(err)
				return
			}
			_ = gzipWriter.Close()
			request, err := http.NewRequest(http.MethodPost, endpoint, &buf)
			if err != nil {
				log.Info(err)
			}

			request.Header.Set("content-type", "application/json")
			request.Header.Set("Content-Encoding", "gzip")
			request.Header.Set("Accept-Encoding", "gzip")

			response, err := client.Do(request)
			if err != nil {
				log.Info(err)
				continue
			}
			defer response.Body.Close()
		}
		for i, v := range counter {
			if i == "" {
				return
			}

			metric := models.Metrics{MType: models.Counter, ID: i, Delta: &v}
			data, err := json.Marshal(metric)
			if err != nil {
				log.Info(err)
				return
			}

			var buf bytes.Buffer

			gzipWriter := gzip.NewWriter(&buf)

			_, err = gzipWriter.Write(data)
			if err != nil {
				log.Info(err)
				return
			}

			err = gzipWriter.Flush()
			if err != nil {
				log.Info(err)
				return
			}
			_ = gzipWriter.Close()
			request, err := http.NewRequest(http.MethodPost, endpoint, &buf)
			if err != nil {
				log.Info(err)
				return
			}

			request.Header.Set("content-type", "application/json")
			request.Header.Set("Content-Encoding", "gzip")
			request.Header.Set("Accept-Encoding", "gzip")
			response, err := client.Do(request)
			if err != nil {
				log.Info(err)
				continue
			}
			defer response.Body.Close()
		}
	}
}
