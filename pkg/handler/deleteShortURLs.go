package handler

import (
	"encoding/json"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/utils"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) DeleteShortURLs(w http.ResponseWriter, r *http.Request) {
	var URLsList []string

	ctx := r.Context()
	UserID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		logger.Log.Error("Cannot get UserID from context")
	}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&URLsList); err != nil {
		logger.Log.Error("Error decode request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.services.DelShortURLs(UserID, URLsList)

	w.WriteHeader(http.StatusAccepted)
}
