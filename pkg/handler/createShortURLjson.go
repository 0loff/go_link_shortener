package handler

import (
	"encoding/json"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/models"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) CreateShortURLjson(w http.ResponseWriter, r *http.Request) {
	var origURL models.CreateURLRequestPayload

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&origURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("Cannot decode request JSON body", zap.Error(err))
		return
	}

	if string(origURL.URL) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := models.CreateURLResponsePayload{
		Result: h.services.ShortURLHost + "/" + h.services.ShortURLResolver(string(origURL.URL)),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("Error encoding response ", zap.Error(err))
		return
	}
}
