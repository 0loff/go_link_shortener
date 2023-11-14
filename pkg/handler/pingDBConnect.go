package handler

import (
	"context"
	"go_link_shortener/internal/logger"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) PingConnect(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	err := h.services.Repo.PingConnect(ctx)
	if err != nil {
		logger.Log.Error("Error database connection", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
