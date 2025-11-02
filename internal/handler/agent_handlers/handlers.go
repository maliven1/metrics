package agenthandlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/maliven1/metrics/internal/config"
	models "github.com/maliven1/metrics/internal/model"
	"go.uber.org/zap"
)

type Agent interface {
	GetMetrics() (map[string]float64, map[string]int64)
	CollectMetrics(m *sync.Mutex)
	MakeHash(value string, key string) string
}

type SendClient struct {
	AddHandler Agent
	cfg        *config.AgentConfig
}

func NewSendClient(s Agent, cfg *config.AgentConfig) *SendClient {
	return &SendClient{AddHandler: s, cfg: cfg}
}

// Semaphore структура семафора
type Semaphore struct {
	semaCh chan struct{}
}

// NewSemaphore создает семафор с буферизованным каналом емкостью maxReq
func NewSemaphore(maxReq int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

// когда горутина запускается, отправляем пустую структуру в канал semaCh
func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

// когда горутина завершается, из канала semaCh убирается пустая структура
func (s *Semaphore) Release() {
	<-s.semaCh
}

func (s SendClient) SendClientMetrics() {
	endpoint := "http://" + s.cfg.Address + "/update/"
	client := &http.Client{}
	var m sync.Mutex
	go s.AddHandler.CollectMetrics(&m)
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
			hash := s.AddHandler.MakeHash("", s.cfg.Key)
			request.Header.Set("HashSHA256", hash)
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
			hash := s.AddHandler.MakeHash("", s.cfg.Key)
			request.Header.Set("HashSHA256", hash)
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

func (s SendClient) SendClientBatchMetrics(log *zap.SugaredLogger, wg *sync.WaitGroup) {
	var delay = time.Second           // Начальная задержка
	const increment = 2 * time.Second // Увеличение задержки на 2 секунды после каждой попытки

	endpoint := "http://" + s.cfg.Address + "/updates/"
	log.Info("start agent on endpoint: ", endpoint)
	client := &http.Client{}
	var m sync.Mutex
	go s.AddHandler.CollectMetrics(&m)
	semaphore := NewSemaphore(s.cfg.RateLimit)

	for {
		semaphore.Acquire()
		go func() {
			defer semaphore.Release()

			time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
			err := retry.Do(func() error {
				gauge, counter := s.AddHandler.GetMetrics()

				var metrics []models.Metrics

				for i, v := range gauge {
					if i == "" {
						return nil
					}
					metric := models.Metrics{MType: models.Gauge, ID: i, Value: &v}
					metrics = append(metrics, metric)
				}

				for i, v := range counter {
					if i == "" {
						return nil
					}
					metric := models.Metrics{MType: models.Counter, ID: i, Delta: &v}
					metrics = append(metrics, metric)
				}

				if len(metrics) > 0 {
					data, err := json.Marshal(metrics)
					if err != nil {
						log.Error("Failed to marshal metrics batch: ", err)
						return nil
					}

					var buf bytes.Buffer
					gzipWriter := gzip.NewWriter(&buf)

					_, err = gzipWriter.Write(data)
					if err != nil {
						log.Error("Failed to write to gzip writer: ", err)
						return nil
					}
					err = gzipWriter.Flush()
					if err != nil {
						log.Error("Failed to flush gzip writer: ", err)
						return nil
					}
					err = gzipWriter.Close()
					if err != nil {
						log.Error("Failed to close gzip writer: ", err)
						return nil
					}

					request, err := http.NewRequest(http.MethodPost, endpoint, &buf)
					if err != nil {
						log.Error("Failed to create request: ", err)
						return nil
					}

					hash := s.AddHandler.MakeHash(string(data), s.cfg.Key)
					request.Header.Set("HashSHA256", hash)
					request.Header.Set("content-type", "application/json")
					request.Header.Set("Content-Encoding", "gzip")
					request.Header.Set("Accept-Encoding", "gzip")

					response, err := client.Do(request)
					if err != nil {
						log.Info(err)
						return err
					}
					defer response.Body.Close()

				}
				return nil
			}, retry.Attempts(3), retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
				if n > 0 {
					delay += increment
				}
				return delay
			}))
			if err != nil {
				return
			}

		}()
	}

}

func (s SendClient) SendClientJSONMetrics(log *zap.SugaredLogger, wg *sync.WaitGroup) {
	var delay = time.Second           // Начальная задержка
	const increment = 2 * time.Second // Увеличение задержки на 2 секунды после каждой попытки

	endpoint := "http://" + s.cfg.Address + "/update/"
	log.Info("start agent on endpoint: ", endpoint)
	client := &http.Client{}

	var m sync.Mutex

	defer wg.Done()
	semaphore := NewSemaphore(s.cfg.RateLimit)

	go s.AddHandler.CollectMetrics(&m)

	for {
		semaphore.Acquire()

		go func() {
			defer semaphore.Release()
			time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
			err := retry.Do(func() error {
				gauge, counter := s.AddHandler.GetMetrics()
				for i, v := range gauge {
					if i == "" {
						return nil
					}

					metric := models.Metrics{MType: models.Gauge, ID: i, Value: &v}
					data, err := json.Marshal(metric)
					if err != nil {
						log.Error(err)
						return nil
					}

					var buf bytes.Buffer
					gzipWriter := gzip.NewWriter(&buf)

					_, err = gzipWriter.Write(data)
					if err != nil {
						log.Error(err)
						return nil
					}
					err = gzipWriter.Flush()
					if err != nil {
						log.Info(err)
						return nil
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

					hash := s.AddHandler.MakeHash(string(data), s.cfg.Key)
					request.Header.Set("HashSHA256", hash)

					request.Header.Set("content-type", "application/json")
					request.Header.Set("Content-Encoding", "gzip")
					request.Header.Set("Accept-Encoding", "gzip")
					response, err := client.Do(request)
					if err != nil {
						log.Info(err)
						return err
					}
					defer response.Body.Close()

				}

				for i, v := range counter {

					if i == "" {
						break
					}

					metric := models.Metrics{MType: models.Counter, ID: i, Delta: &v}
					data, err := json.Marshal(metric)
					if err != nil {
						log.Info(err)

						return err
					}

					var buf bytes.Buffer

					gzipWriter := gzip.NewWriter(&buf)

					_, err = gzipWriter.Write(data)
					if err != nil {
						log.Info(err)

						return err
					}

					err = gzipWriter.Flush()
					if err != nil {
						log.Info(err)

						return err
					}
					_ = gzipWriter.Close()
					request, err := http.NewRequest(http.MethodPost, endpoint, &buf)
					if err != nil {
						log.Info(err)

						return err
					}

					hash := s.AddHandler.MakeHash(string(data), s.cfg.Key)
					request.Header.Set("HashSHA256", hash)

					request.Header.Set("content-type", "application/json")
					request.Header.Set("Content-Encoding", "gzip")
					request.Header.Set("Accept-Encoding", "gzip")

					response, err := client.Do(request)
					if err != nil {
						log.Info(err)
						return err
					}
					defer response.Body.Close()
				}
				return nil
			}, retry.Attempts(3), retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
				if n > 0 {
					delay += increment
				}
				return delay
			}))
			if err != nil {
				return
			}
		}()
	}

}
