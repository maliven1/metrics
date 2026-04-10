package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"

	"github.com/maliven1/metrics/internal/config"
	crypto "github.com/maliven1/metrics/internal/crypto"
	"go.uber.org/zap"
)

func DecryptedMessage(cfg config.ServerConfig, log *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ow := w
			secretKey, err := crypto.ReadKeys(cfg)
			if err != nil {
				if err.Error() == "private key not found" {
					h.ServeHTTP(ow, r)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				log.Error(http.StatusInternalServerError, err)
				return
			}

			// Читаем тело запроса
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Info(http.StatusBadRequest, err)
				return
			}
			defer r.Body.Close()

			// Восстанавливаем тело запроса на случай, если не будем дешифровать
			shouldRestore := true
			defer func() {
				if shouldRestore {
					r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				}
			}()

			// Проверяем, что данные не пустые
			if len(bodyBytes) == 0 {
				h.ServeHTTP(ow, r)
				return
			}

			// Проверяем, что данные имеют правильный размер для RSA
			keySize := secretKey.Size()
			if len(bodyBytes) != keySize {
				// Проверяем, может быть это уже JSON (начинается с {)
				// Если это JSON, пропускаем дешифровку
				h.ServeHTTP(ow, r)
				return
			}

			decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, secretKey, bodyBytes)
			if err != nil {
				log.Errorf("RSA decryption failed: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				log.Error(http.StatusInternalServerError, err)
				return
			}
			// Заменяем тело запроса на расшифрованные данные
			r.Body = io.NopCloser(bytes.NewReader(decryptedMessage))
			// Отменяем восстановление оригинального тела
			shouldRestore = false

			h.ServeHTTP(ow, r)
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
