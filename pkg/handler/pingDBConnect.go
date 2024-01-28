package handler

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/0loff/go_link_shortener/internal/logger"
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
