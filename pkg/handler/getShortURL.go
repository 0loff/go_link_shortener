package handler

import (
	"errors"
	"go_link_shortener/internal/logger"
	"go_link_shortener/pkg/repository"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) GetShortURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	link, err := h.services.GetShortURL(ctx, r.URL.Path[1:])
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrURLGone):
			logger.Log.Error("Cannot get URL", zap.Error(err))
			w.WriteHeader(http.StatusGone)
			return

		default:
			logger.Log.Error("Not found", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
