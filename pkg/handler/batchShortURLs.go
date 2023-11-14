package handler

import (
	"context"
	"encoding/json"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/models"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) BatchShortURLs(w http.ResponseWriter, r *http.Request) {
	entries := []models.BatchURLRequestEntry{}

	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("Error parsing request body", zap.Error(err))
		return
	}

	err = json.Unmarshal(body, &entries)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("Error parsing JSON", zap.Error(err))
		return
	}

	resp := h.services.SetBatchShortURLs(ctx, entries)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("Error encoding response ", zap.Error(err))
		return
	}
}
