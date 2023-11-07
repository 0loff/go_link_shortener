package handler

import (
	"net/http"
)

func (h *Handler) PingDBConnect(w http.ResponseWriter, r *http.Request) {
	err := h.services.Repo.PingDBConnect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
