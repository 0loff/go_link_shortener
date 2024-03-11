package handler

import (
	"encoding/json"
	"net/http"

	"github.com/0loff/go_link_shortener/internal/logger"
	"go.uber.org/zap"
)

// GetStats handler for reseived statistics about amount of short urls and active users
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	if err := enc.Encode(h.services.GetStatistics()); err != nil {
		logger.Log.Error("Error encoding response ", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
