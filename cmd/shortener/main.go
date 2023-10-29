package main

import (
	"go_link_shortener/internal/handlers"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/storage"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func gzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")

		if supportGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {

			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}
		h.ServeHTTP(ow, r)
	})
}

func CustomRouter() chi.Router {
	r := chi.NewRouter()

	return r.Route("/", func(r chi.Router) {
		r.Use(gzipMiddleware)
		r.Use(logger.RequestLogger)

		r.Post("/", http.HandlerFunc(handlers.CreateLinkHandler))
		r.Get("/{id}", http.HandlerFunc(handlers.GetLinkHandler))
		r.Post("/api/shorten", http.HandlerFunc(handlers.CreateLinkJSONHandler))
	})
}

func main() {
	NewConfigBuilder()
	storage.LinkStorageInit()
	storage.Store.SetStorageFile(config.StorageFile)
	storage.Store.SetShortLinkHost(config.ShortLinkHost)
	storage.Store.LinkStorageRecover()

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
