package handler

import (
	"encoding/json"
	"fmt"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/models"
	"io"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) BatchShortURLs(w http.ResponseWriter, r *http.Request) {
	entries := []models.BatchURLRequestEntry{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(body, &entries)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp := h.services.SetBatchShortURLs(entries)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("Error encoding response ", zap.Error(err))
		return
	}
}
