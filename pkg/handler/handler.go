package handler

import (
	"go_link_shortener/internal/compressor"
	"go_link_shortener/internal/logger"
	"go_link_shortener/pkg/service"
	"net/http"

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
		r.Use(compressor.GzipCompressor)
		r.Use(logger.RequestLogger)

		r.Post("/", http.HandlerFunc(h.CreateShortURL))
		r.Post("/api/shorten", http.HandlerFunc(h.CreateShortURLjson))
		r.Get("/{id}", http.HandlerFunc(h.GetShortURL))
		r.Get("/ping", http.HandlerFunc(h.PingConnect))
	})
}
