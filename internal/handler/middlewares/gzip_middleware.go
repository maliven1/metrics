// Package middlewares
package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func GzipMiddleware(log *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
			// который будем передавать следующей функции
			ow := w

			// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
				cw := NewCompressWriter(w)
				// меняем оригинальный http.ResponseWriter на новый
				ow = cw
				// не забываем отправить клиенту все сжатые данные после завершения middleware
				defer cw.Close()
			}

			// проверяем, что клиент отправил серверу сжатые данные в формате gzip
			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				// Читаем всё тело запроса в буфер, чтобы проверить, является ли оно gzip
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					log.Warnf("Failed to read request body: %v", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				// Восстанавливаем тело запроса из буфера
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

				// Проверяем, является ли тело валидным gzip (проверяем magic number)
				isGzip := len(bodyBytes) >= 2 && bodyBytes[0] == 0x1f && bodyBytes[1] == 0x8b

				if isGzip {
					// Создаём gzip reader из буфера
					cr, err := NewCompressReader(io.NopCloser(bytes.NewReader(bodyBytes)))
					if err != nil {
						log.Warnf("Failed to create gzip reader for valid gzip data: %v. Skipping decompression.", err)
						// Оставляем тело как есть (raw gzip данные)
						// Хендлер должен будет разобраться с этим
					} else {
						// меняем тело запроса на новое
						r.Body = cr
						defer cr.Close()
					}
				} else {
					// Safely get first 8 bytes or fewer
					firstBytes := bodyBytes
					if len(firstBytes) > 8 {
						firstBytes = firstBytes[:8]
					}
					// Оставляем тело как есть (вероятно, это уже распакованные данные)
				}
			}

			// передаём управление хендлеру
			h.ServeHTTP(ow, r)
		})
	}
}

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
