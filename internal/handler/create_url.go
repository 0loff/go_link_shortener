package handler

import (
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/utils"
)

// Обработчик создания и сохранения сокращенного URL пользователя, при передаче в формате text/plain тела запроса
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	UserID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		logger.Log.Error("Cannot get UserID from context")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("Error parsing request body", zap.Error(err))
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statusHeader := http.StatusCreated
	shortURL, err := h.services.CreateShortURL(ctx, UserID, string(body))
	if err != nil {
		statusHeader = http.StatusConflict

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")

		if supportGzip {
			w.Header().Set("Content-Encoding", "gzip")
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusHeader)
	w.Write([]byte(h.services.ShortURLHost + "/" + shortURL))
}
