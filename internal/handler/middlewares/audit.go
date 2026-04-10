// Package middlewares
package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/maliven1/metrics/internal/config"
	models "github.com/maliven1/metrics/internal/model"
	"go.uber.org/zap"
)

type AuditEvent struct {
	Timestamp int64    `json:"ts"`
	Metrics   []string `json:"metrics"`
	IPAddress string   `json:"ip_address"`
}

type AuditReceiver interface {
	Notify(event *AuditEvent) error
}

type FileAuditReceiver struct {
	FilePath string
}

func (f *FileAuditReceiver) Notify(event *AuditEvent) error {
	file, err := os.OpenFile(f.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = file.Write(append(data, '\n'))
	return err
}

type URLAuditReceiver struct {
	URL string
}

func (u *URLAuditReceiver) Notify(event *AuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	resp, err := http.Post(u.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send audit log, status: %s", resp.Status)
	}
	return nil
}

// AuditMiddleware - извлекает метрики и передает их в аудит
func AuditMiddleware(log *zap.SugaredLogger, cfg config.ServerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if cfg.AuditFilePath == "" && cfg.AuditURL == "" {

				next.ServeHTTP(w, r)
				return
			}

			auditReceivers := make([]AuditReceiver, 0, 2)
			if cfg.AuditFilePath != "" {
				auditReceivers = append(auditReceivers, &FileAuditReceiver{FilePath: cfg.AuditFilePath})
			}
			if cfg.AuditURL != "" {
				auditReceivers = append(auditReceivers, &URLAuditReceiver{URL: cfg.AuditURL})
			}

			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				log.Info("Error reading request body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(buf.Bytes()))

			metrics := []models.Metrics{}
			metric := models.Metrics{}
			if err := json.Unmarshal(buf.Bytes(), &metrics); err != nil {
				if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				metrics = append(metrics, metric)
			}

			log.Info("Extracted metrics: %v", metrics)

			var auditMetrics []string
			for _, metric := range metrics {
				auditMetrics = append(auditMetrics, metric.ID)
			}

			if len(auditMetrics) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			event := &AuditEvent{
				Timestamp: time.Now().Unix(),
				Metrics:   auditMetrics,
				IPAddress: r.RemoteAddr,
			}

			go func() {
				for _, receiver := range auditReceivers {
					if err := receiver.Notify(event); err != nil {
						log.Warnf("Failed to send audit event: %v", err)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
