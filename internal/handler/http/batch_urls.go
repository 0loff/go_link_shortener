package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/models"
	"github.com/0loff/go_link_shortener/internal/utils"
)

// BatchShortURLs - Handler for batch insertion of short URLs by user
func (h *Handler) BatchShortURLs(w http.ResponseWriter, r *http.Request) {
	entries := []models.BatchURLRequestEntry{}

	ctx := r.Context()
	UserID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		logger.Log.Error("Cannot get UserID from context")
	}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&entries); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("Error json decode request body", zap.Error(err))
		return
	}

	resp := h.services.SetBatchShortURLs(ctx, UserID, entries)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("Error encoding response ", zap.Error(err))
		return
	}
}
