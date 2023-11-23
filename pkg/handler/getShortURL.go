package handler

import (
	"net/http"
)

func (h *Handler) GetShortURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	link := h.services.GetShortURL(ctx, r.URL.Path[1:])

	if link == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
