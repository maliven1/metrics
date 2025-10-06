package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maliven1/metrics/internal/handler/middlewares"
	serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers"
	"github.com/maliven1/metrics/internal/logger"
	"go.uber.org/zap"
)

func NewRouter(r *chi.Mux, handler *serverhandlers.AddHandler, log *zap.SugaredLogger) {
	r.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return logger.WithLogging(h, log)
		})
		r.Group(func(r chi.Router) {

			r.Use(middlewares.GzipMiddleware(log))

			r.Get(`/`, handler.GetAllMetricsHandler())
			r.Post(`/value/`, handler.GetBodyMetricHandler(log))
			r.Post(`/update/`, handler.PostBodyHandler(log))
		})
		r.Post(`/update/*`, handler.PostURLHandler())
		r.Get(`/value/*`, handler.GetMetricHandler())

	})

}
