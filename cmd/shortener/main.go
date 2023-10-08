package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var linkStorage = map[string]string{}

func createLink(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	linkStorage[base64.RawStdEncoding.EncodeToString(body)] = string(body)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + config.ShortLinkHost + "/" + base64.RawStdEncoding.EncodeToString(body)))
}

func getLink(w http.ResponseWriter, r *http.Request) {

	link, ok := linkStorage[r.URL.Path[1:]]

	if ok {
		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func CustomRouter() chi.Router {
	r := chi.NewRouter()

	return r.Route("/", func(r chi.Router) {
		r.Post("/", createLink)
		r.Get("/{id}", getLink)
	})
}

func notAllowedRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	NewConfigBuilder()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Running server on", config.Host)
	return http.ListenAndServe(config.Host, CustomRouter())
}
