package handler

import "net/http"

func (h *Handler) GetShortURL(w http.ResponseWriter, r *http.Request) {
	link := h.services.GetShortURL(r.URL.Path[1:])

	if link != "" {
		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
