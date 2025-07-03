package book

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Get("/", h.GetAllBooks)
	r.Get("/{id}", h.GetBookByID)
	r.Post("/", h.AddBook)
	r.Put("/", h.UpdateBook)
	r.Delete("/{id}", h.DeleteBook)

	r.Get("/search/google", h.SearchGoogleBooks)
}
