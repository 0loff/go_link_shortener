package main

import (
	"go_link_shortener/internal/handlers"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CustomRouter() chi.Router {
	r := chi.NewRouter()

	return r.Route("/", func(r chi.Router) {
		r.Post("/", logger.RequestLogger(http.HandlerFunc(handlers.CreateLinkHandler)))
		r.Get("/{id}", logger.RequestLogger(http.HandlerFunc(handlers.GetLinkHandler)))
	})
}

func main() {
	NewConfigBuilder()
	storage.LinkStorageInit()
	storage.Store.SetShortLinkHost(config.ShortLinkHost)

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize(config.LogLevel); err != nil {
		return err
	}

	logger.Sugar.Infoln("Host", config.Host)

	return http.ListenAndServe(config.Host, CustomRouter())
}
