package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/utils"
)

func (h *Handler) GetShortURLs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	UserID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		logger.Log.Error("Cannot get UserID from context")
	}

	UserStatus, ok := utils.GetUserStatusFromContext(ctx)
	if ok && UserStatus == "new_user" {
		w.WriteHeader(http.StatusUnauthorized)
		return
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
