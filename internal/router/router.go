package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maliven1/metrics/internal/config"
	"github.com/maliven1/metrics/internal/handler/middlewares"
	serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers"
	"github.com/maliven1/metrics/internal/logger"
	"go.uber.org/zap"
)

func NewRouter(r *chi.Mux, handler *serverhandlers.Handler, log *zap.SugaredLogger, cfg config.ServerConfig) {
	r.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return logger.WithLogging(h, log)
		})
		r.Group(func(r chi.Router) {

			r.Use(middlewares.GzipMiddleware(log))

			r.Get(`/`, handler.GetAllMetricsHandler())
			r.Post(`/value/`, handler.GetBodyMetricHandler(log))
			r.Group(func(r chi.Router) {
				r.Use(middlewares.HashMiddleware(log, cfg))
				r.Post(`/update/`, handler.PostBodyHandler(log))
				r.Post(`/updates/`, handler.PostMetricsHandler(log))
				r.Post(`/update/*`, handler.PostURLHandler())
			})
		})

		r.Get(`/value/*`, handler.GetMetricHandler())
		r.Get(`/ping`, handler.PingHandler(log))

	})

}
