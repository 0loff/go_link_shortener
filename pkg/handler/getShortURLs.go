package handler

import (
	"encoding/json"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/utils"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) GetShortURLs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	UserID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		logger.Log.Error("Cannot get UserID from context")
	}

	resp := h.services.GetShortURLs(ctx, UserID)

	w.Header().Set("Content-Type", "application/json")
	if len(resp) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("Error encoding response ", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
