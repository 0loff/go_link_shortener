package handlers

import (
	"go_link_shortener/internal/base62"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/storage"
	"io"
	"log"
	"net/http"

	"go.uber.org/zap"
)

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

	var shortURL string
	existedShortURL := storage.Store.FindByLink(string(body))
	if existedShortURL != "" {
		shortURL = existedShortURL
	} else {
		shortURL = base62.EncodedString()
		storage.Store.SetShortLink(shortURL, string(body))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(storage.Store.GetShortLinkHost() + "/" + shortURL))
}

func GetLinkHandler(w http.ResponseWriter, r *http.Request) {

	link, ok := storage.Store.GetShortLink(r.URL.Path[1:])

	logger.Log.Debug("Got request method ", zap.String("method", r.Method))

	if ok {
		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
