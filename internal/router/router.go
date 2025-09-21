package router

import (
	"github.com/go-chi/chi"
	serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers"
)

func New(h serverhandlers.AddHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Post(`/update/*`, serverhandlers.PostHandler(h))
	router.Get(`/value/*`, serverhandlers.GetMetricHandler(h))
	router.Get(`/`, serverhandlers.GetAllMetricsHandler(h))

	return router

}
