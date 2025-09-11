package handler

import (
	"net/http"
	"strings"

	models "github.com/maliven1/metrics/internal/model"
	"github.com/maliven1/metrics/internal/service"
)

func PostHandler(memStorage models.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		if r.Method == http.MethodPost {
			pathSplit := strings.Split(r.URL.Path, "/")
			status := service.CheckPath(pathSplit, memStorage)
			w.WriteHeader(status)
		} else {

			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
