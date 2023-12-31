package handler

import (
	"context"
	"encoding/json"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/models"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func (h *Handler) CreateShortURLjson(w http.ResponseWriter, r *http.Request) {
	var origURL models.CreateURLRequestPayload

	ctx := context.Background()

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&origURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("Cannot decode request JSON body", zap.Error(err))
		return
	}

	if string(origURL.URL) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statusHeader := http.StatusCreated
	shortURL, err := h.services.CreateShortURL(ctx, origURL.URL)
	if err != nil {
		statusHeader = http.StatusConflict

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")

		if supportGzip {
			w.Header().Set("Content-Encoding", "gzip")
		}
	}

	resp := models.CreateURLResponsePayload{
		Result: h.services.ShortURLHost + "/" + shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusHeader)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("Error encoding response ", zap.Error(err))
		return
	}
}
