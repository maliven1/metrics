// Package agenthandlers
package agenthandlers

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/maliven1/metrics/internal/agent"
	"github.com/maliven1/metrics/internal/config"
	crypto "github.com/maliven1/metrics/internal/crypto"
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
	Client     *http.Client
}

func NewSendClient(s Agent, cfg *config.AgentConfig) *SendClient {
	return &SendClient{
		AddHandler: s,
		cfg:        cfg,
		Client:     &http.Client{},
	}
}

// NewSendClientWithHTTPClient creates a new SendClient with a custom HTTP client
func NewSendClientWithHTTPClient(s Agent, cfg *config.AgentConfig, client *http.Client) *SendClient {
	return &SendClient{
		AddHandler: s,
		cfg:        cfg,
		Client:     client,
	}
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

			hash := crypto.MakeHash("", s.cfg.Key)
			request.Header.Set("HashSHA256", hash)
			request.Header.Set("Content-Type", "Content-Type: text/plain")

			response, err := s.Client.Do(request)
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
			hash := crypto.MakeHash("", s.cfg.Key)
			request.Header.Set("HashSHA256", hash)
			request.Header.Set("Content-Type", "Content-Type: text/plain")

			response, err := s.Client.Do(request)
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

	defer wg.Done()
	endpoint := "http://" + s.cfg.Address + "/updates/"
	log.Info("start agent on endpoint: ", endpoint)

	certificate := agent.ReadKey(s.cfg)
	go s.AddHandler.CollectMetrics()
	var w sync.WaitGroup
	w.Add(s.cfg.RateLimit)
	for i := 0; i < s.cfg.RateLimit; i++ {
		go func() {
			defer w.Done()
			for {
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
						encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, certificate.PublicKey.(*rsa.PublicKey), buf.Bytes())
						if err != nil {
							log.Fatal(err)
						}
						request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(encryptedMessage))
						if err != nil {
							log.Error("Failed to create request: ", err)
							return nil
						}

						hash := crypto.MakeHash(string(data), s.cfg.Key)
						request.Header.Set("HashSHA256", hash)
						request.Header.Set("content-type", "application/json")
						request.Header.Set("Content-Encoding", "gzip")
						request.Header.Set("Accept-Encoding", "gzip")

						response, err := s.Client.Do(request)
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

			}
		}()
	}
	w.Wait()
}

func (s SendClient) SendClientJSONMetrics(log *zap.SugaredLogger, wg *sync.WaitGroup) {
	var delay = time.Second           // Начальная задержка
	const increment = 2 * time.Second // Увеличение задержки на 2 секунды после каждой попытки
	certificate := agent.ReadKey(s.cfg)
	endpoint := "http://" + s.cfg.Address + "/update/"
	log.Info("start agent on endpoint: ", endpoint)

	defer wg.Done()
	var w sync.WaitGroup
	go s.AddHandler.CollectMetrics()
	w.Add(s.cfg.RateLimit)
	for i := 0; i < s.cfg.RateLimit; i++ {
		go func() {
			defer w.Done()
			for {

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
						encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, certificate.PublicKey.(*rsa.PublicKey), buf.Bytes())
						if err != nil {
							log.Fatal(err)
						}
						request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(encryptedMessage))
						if err != nil {
							log.Info(err)
						}

						hash := crypto.MakeHash(string(data), s.cfg.Key)
						request.Header.Set("HashSHA256", hash)

						request.Header.Set("content-type", "application/json")
						request.Header.Set("Content-Encoding", "gzip")
						request.Header.Set("Accept-Encoding", "gzip")
						response, err := s.Client.Do(request)
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

						hash := crypto.MakeHash(string(data), s.cfg.Key)
						request.Header.Set("HashSHA256", hash)

						request.Header.Set("content-type", "application/json")
						request.Header.Set("Content-Encoding", "gzip")
						request.Header.Set("Accept-Encoding", "gzip")

						response, err := s.Client.Do(request)
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

			}
		}()

	}
	w.Wait()
}
