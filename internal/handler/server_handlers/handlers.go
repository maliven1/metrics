package serverhandlers

import (
	"net/http"
	"strings"
)

type Service interface {
	CheckPath(pathSplit []string) int
}

type AddHandler struct {
	AddHandler Service
}

func NewAddHandler(s Service) *AddHandler {
	return &AddHandler{AddHandler: s}
}

func (h AddHandler) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if r.Method == http.MethodPost {
			pathSplit := strings.Split(r.URL.Path, "/")
			status := h.AddHandler.CheckPath(pathSplit)
			w.WriteHeader(status)
		} else {

			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
