package handler

import (
	"go_link_shortener/internal/auth"
	"go_link_shortener/internal/compressor"
	"go_link_shortener/internal/logger"
	"go_link_shortener/pkg/service"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() chi.Router {
	r := chi.NewRouter()

	return r.Route("/", func(r chi.Router) {
		r.Group(func(router chi.Router) {
			router.Use(compressor.GzipCompressor)
			router.Use(logger.RequestLogger)
			router.Use(auth.UserAuth)

			router.Get("/{id}", http.HandlerFunc(h.GetShortURL))
			router.Get("/ping", http.HandlerFunc(h.PingConnect))
			router.Get("/api/user/urls", http.HandlerFunc(h.GetShortURLs))

			router.Post("/", http.HandlerFunc(h.CreateShortURL))
			router.Post("/api/shorten", http.HandlerFunc(h.CreateShortURLjson))
			router.Post("/api/shorten/batch", http.HandlerFunc(h.BatchShortURLs))

			router.Delete("/api/user/urls", http.HandlerFunc(h.DeleteShortURLs))
		})

		r.Get("/debug/pprof/", pprof.Index)
		r.Get("/debug/pprof/cmdline", pprof.Cmdline)
		r.Get("/debug/pprof/profile", pprof.Profile)
		r.Get("/debug/pprof/symbol", pprof.Symbol)
		r.Get("/debug/pprof/trace", pprof.Trace)
		r.Get("/debug/pprof/{cmd}", pprof.Index)
	})
}
