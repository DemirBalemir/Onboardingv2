package author

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Post("/", h.RegisterAuthor)
	r.Get("/{id}", h.GetAuthorByID)
}
