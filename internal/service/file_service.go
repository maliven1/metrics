package service

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/maliven1/metrics/internal/config"
	models "github.com/maliven1/metrics/internal/model"
	"go.uber.org/zap"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteMetric(metric *models.Metrics) error {
	return p.encoder.Encode(&metric)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadMetric() (*models.Metrics, error) {
	metric := &models.Metrics{}
	if err := c.decoder.Decode(&metric); err != nil {
		return nil, fmt.Errorf("failed to decode metric: %w", err)
	}

	return metric, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

func (s Service) ReadFileMetrics(cfg config.ServerConfig, log *zap.SugaredLogger) {

	if !cfg.Restore {
		return
	}

	Consumer, err := NewConsumer(cfg.FileStoragePath)
	if err != nil {
		log.Error(err)
	}
	defer Consumer.Close()
	for {
		metrics, err := Consumer.ReadMetric()
		if err != nil {

			log.Error(err)
			return
		}
		s.AddStructMetric(*metrics)
	}
}

func (s Service) WriteFileMetrics(cfg config.ServerConfig, log *zap.SugaredLogger) {

	metrics := &models.Metrics{}

	for {

		time.Sleep(time.Duration(cfg.StoreInterval) * time.Second)
		Producer, err := NewProducer(cfg.FileStoragePath)
		if err != nil {
			log.Error(err)
		}
		counter, gauge := s.GetAllMetrics()
		for i, v := range counter {
			if i != "" {
				metrics.MType = models.Counter
				metrics.ID = i
				metrics.Delta = &v
				err = Producer.WriteMetric(metrics)
				if err != nil {
					log.Error(err)
					continue
				}
			}

		}
		for i, v := range gauge {
			if i != "" {
				metrics.MType = models.Gauge
				metrics.ID = i
				metrics.Value = &v
				err = Producer.WriteMetric(metrics)
				if err != nil {
					log.Error(err)
					continue
				}
			}
		}
		Producer.Close()
	}

}

func (s Service) InitFile(cfg config.ServerConfig, log *zap.SugaredLogger) {
	s.ReadFileMetrics(cfg, log)
	s.WriteFileMetrics(cfg, log)
}
