package agenthandlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
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

func (s SendClient) SendClientBatchMetrics(log *zap.SugaredLogger, wg sync.WaitGroup) {
	const maxRetries = 3
	const interval = 2
	var latsError error
	endpoint := "http://" + s.cfg.Address + "/updates/"
	log.Info("start agent on endpoint: ", endpoint)
	client := &http.Client{}
	go s.AddHandler.CollectMetrics()
	for {
		if latsError != nil {
			wg.Done()
			return
		}
		LinearBackoff := 1
		time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
		for attempt := 0; attempt <= maxRetries; {
			gauge, counter := s.AddHandler.GetMetrics()

			// Collect all metrics in a batch
			var metrics []models.Metrics

			// Add gauge metrics to batch
			for i, v := range gauge {
				if i == "" {
					continue
				}
				metric := models.Metrics{MType: models.Gauge, ID: i, Value: &v}
				metrics = append(metrics, metric)
			}

			// Add counter metrics to batch
			for i, v := range counter {
				if i == "" {
					continue
				}
				metric := models.Metrics{MType: models.Counter, ID: i, Delta: &v}
				metrics = append(metrics, metric)
			}

			// Send batch of metrics
			if len(metrics) > 0 {
				data, err := json.Marshal(metrics)
				if err != nil {
					log.Error("Failed to marshal metrics batch: ", err)
					continue
				}

				var buf bytes.Buffer
				gzipWriter := gzip.NewWriter(&buf)

				_, err = gzipWriter.Write(data)
				if err != nil {
					log.Error("Failed to write to gzip writer: ", err)
					continue
				}
				err = gzipWriter.Flush()
				if err != nil {
					log.Error("Failed to flush gzip writer: ", err)
					continue
				}
				err = gzipWriter.Close()
				if err != nil {
					log.Error("Failed to close gzip writer: ", err)
					continue
				}

				request, err := http.NewRequest(http.MethodPost, endpoint, &buf)
				if err != nil {
					log.Error("Failed to create request: ", err)
					continue
				}

				request.Header.Set("content-type", "application/json")
				request.Header.Set("Content-Encoding", "gzip")
				request.Header.Set("Accept-Encoding", "gzip")

				response, err := client.Do(request)
				if err != nil {
					if errors.Is(err, err.(net.Error)) {
						log.Info(err)
						LinearBackoff += interval
						attempt++
						latsError = err
						break
					} else {
						log.Info(err)
						break
					}

				} else {
					latsError = nil
				}
				defer response.Body.Close()
			}
			if latsError == nil {
				break
			}
		}
	}
}

func (s SendClient) SendClientJSONMetrics(log *zap.SugaredLogger, wg sync.WaitGroup) {
	const maxRetries = 3
	const interval = 2
	var latsError error
	endpoint := "http://" + s.cfg.Address + "/update/"
	log.Info("start agent on endpoint: ", endpoint)
	client := &http.Client{}
	go s.AddHandler.CollectMetrics()
	for {
		if latsError != nil {
			wg.Done()
			return
		}
		LinearBackoff := 1
		time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
		for attempt := 0; attempt <= maxRetries; {

			time.Sleep(time.Duration(LinearBackoff) * time.Second)
			gauge, counter := s.AddHandler.GetMetrics()

			for i, v := range gauge {
				if i == "" {
					wg.Done()
					return
				}

				metric := models.Metrics{MType: models.Gauge, ID: i, Value: &v}
				data, err := json.Marshal(metric)
				if err != nil {
					log.Error(err)
					wg.Done()
					return
				}

				var buf bytes.Buffer
				gzipWriter := gzip.NewWriter(&buf)

				_, err = gzipWriter.Write(data)
				if err != nil {
					log.Info(err)
					wg.Done()
					return
				}
				err = gzipWriter.Flush()
				if err != nil {
					log.Info(err)
					wg.Done()
					return
				}
				_ = gzipWriter.Close()

				if buf.Len() == 0 {
					log.Error("Buffer is empty")
					break
				}

				request, err := http.NewRequest(http.MethodPost, endpoint, &buf)
				if err != nil {
					log.Info(err)
				}

				request.Header.Set("content-type", "application/json")
				request.Header.Set("Content-Encoding", "gzip")
				request.Header.Set("Accept-Encoding", "gzip")
				response, err := client.Do(request)
				if err != nil {
					if errors.Is(err, err.(net.Error)) {
						log.Info(err)
						LinearBackoff += interval
						attempt++
						latsError = err
						break
					} else {
						log.Info(err)
						break
					}

				} else {
					latsError = nil
				}
				defer response.Body.Close()
			}

			for i, v := range counter {
				if latsError != nil {
					break
				}
				if i == "" {
					break
				}

				metric := models.Metrics{MType: models.Counter, ID: i, Delta: &v}
				data, err := json.Marshal(metric)
				if err != nil {
					log.Info(err)
					wg.Done()
					return
				}

				var buf bytes.Buffer

				gzipWriter := gzip.NewWriter(&buf)

				_, err = gzipWriter.Write(data)
				if err != nil {
					log.Info(err)
					wg.Done()
					return
				}

				err = gzipWriter.Flush()
				if err != nil {
					log.Info(err)
					wg.Done()
					return
				}
				_ = gzipWriter.Close()
				request, err := http.NewRequest(http.MethodPost, endpoint, &buf)
				if err != nil {
					log.Info(err)
					wg.Done()
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
			if latsError == nil {
				break
			}
		}
	}
}
