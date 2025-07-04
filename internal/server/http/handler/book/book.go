package book

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/demirbalemir/hop/Onboardingv2/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	BookService service.BookService
}

func NewHandler(service service.BookService) *Handler {
	return &Handler{BookService: service}
}

func (h *Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.BookService.GetAllBooks(r.Context())
	if err != nil {
		http.Error(w, "Failed to get books", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	book, err := h.BookService.GetBookByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) AddBook(w http.ResponseWriter, r *http.Request) {
	var book entities.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.BookService.AddBook(r.Context(), &book); err != nil {
		http.Error(w, "Failed to add book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var book entities.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.BookService.UpdateBook(r.Context(), &book); err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	if err := h.BookService.RemoveBook(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SearchGoogleBooks(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Missing title query parameter", http.StatusBadRequest)
		return
	}

	results, err := h.BookService.SearchGoogleBooks(r.Context(), title)
	if err != nil {
		http.Error(w, "Failed to search books", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}
