package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/models"
	"github.com/0loff/go_link_shortener/internal/utils"
)

func (h *Handler) CreateShortURLjson(w http.ResponseWriter, r *http.Request) {
	var origURL models.CreateURLRequestPayload

	ctx := r.Context()
	UserID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		logger.Log.Error("Cannot get UserID from context")
	}

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
	shortURL, err := h.services.CreateShortURL(ctx, UserID, origURL.URL)
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
