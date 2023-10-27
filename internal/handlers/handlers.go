package handlers

import (
	"encoding/json"
	"go_link_shortener/internal/base62"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/models"

	"go_link_shortener/internal/storage"
	"io"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func CreateLinkJSONHandler(w http.ResponseWriter, r *http.Request) {

	var originURL models.CreateURLRequestPayload

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&originURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Debug("Cannot decode request JSON body", zap.Error(err))
		return
	}

	if string(originURL.URL) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := models.CreateURLResponsePayload{
		Result: storage.Store.GetShortLinkHost() + "/" + shortURLResolver(string(originURL.URL)),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("Error encoding response ", zap.Error(err))
		return
	}
}

func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(storage.Store.GetShortLinkHost() + "/" + shortURLResolver(string(body))))
}

func GetLinkHandler(w http.ResponseWriter, r *http.Request) {

	link, ok := storage.Store.GetShortLink(r.URL.Path[1:])

	if ok {
		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func shortURLResolver(url string) string {
	shortURL := storage.Store.FindByLink(url)

	if shortURL != "" {
		return shortURL
	}

	return storage.Store.SetShortLink(base62.EncodedString(), url)
}
