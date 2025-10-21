package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/maliven1/metrics/internal/config"
	"go.uber.org/zap"
)

func HashMiddleware(log *zap.SugaredLogger, cfg config.ServerConfig) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Key == "" {
				h.ServeHTTP(w, r)
				return
			}

			// Read the request body
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Error("buf.ReadFrom(r.Body)", http.StatusBadRequest)
				return
			}

			// Restore the request body for the next handler
			r.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

			// Verify the hash
			serverHash := MakeHash(buf.String(), cfg.Key)
			hashFromHeader := r.Header.Get("HashSHA256")
			if serverHash != hashFromHeader {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Create a response writer that captures the response
			crw := &capturingResponseWriter{
				ResponseWriter: w,
				Buffer:         &bytes.Buffer{},
			}

			// Pass control to the next handler
			h.ServeHTTP(crw, r)

			// Add hash to response headers
			serverHash = MakeHash(crw.Buffer.String(), cfg.Key)
			w.Header().Set("HashSHA256", serverHash)
		})
	}
}

// capturingResponseWriter wraps an http.ResponseWriter to capture the response body
type capturingResponseWriter struct {
	http.ResponseWriter
	Buffer *bytes.Buffer
}

func (cw *capturingResponseWriter) Write(data []byte) (int, error) {
	// Write to both the original ResponseWriter and our buffer
	cw.Buffer.Write(data)
	return cw.ResponseWriter.Write(data)
}

func (cw *capturingResponseWriter) WriteHeader(statusCode int) {
	// Set the Content-Length header to ensure proper handling
	cw.ResponseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", cw.Buffer.Len()))
	cw.ResponseWriter.WriteHeader(statusCode)
}

func MakeHash(value string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	dst := h.Sum(nil)
	return fmt.Sprintf("%x", dst)
}
