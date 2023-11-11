package handler

import (
	"net/http"
)

func (h *Handler) PingConnect(w http.ResponseWriter, r *http.Request) {
	err := h.services.Repo.PingConnect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
