package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func createLink(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var shortURL string
	existedShortURL := storage.FindByLink(string(body))
	if existedShortURL != "" {
		shortURL = existedShortURL
	} else {
		shortURL = shortURLBuilder()
		storage.SetShortLink(shortURL, string(body))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(config.ShortLinkHost + "/" + shortURL))
}

func getLink(w http.ResponseWriter, r *http.Request) {

	link, ok := storage.GetShortLink(r.URL.Path[1:])

	if ok {
		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	return
}

func Base62Encode(id uint64) string {
	length := len(alphabet)
	var encodedBuilder strings.Builder

	encodedBuilder.Grow(10)

	for ; id > 0; id = id / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(id % uint64(length))])
	}

	return encodedBuilder.String()
}

func shortURLBuilder() string {
	return Base62Encode(rand.Uint64())
}

func CustomRouter() chi.Router {
	r := chi.NewRouter()

	return r.Route("/", func(r chi.Router) {
		r.Post("/", createLink)
		r.Get("/{id}", getLink)
	})
}

func main() {
	NewConfigBuilder()
	LinkStorageInit()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Running server on", config.Host)
	return http.ListenAndServe(config.Host, CustomRouter())
}
