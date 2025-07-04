package author

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/demirbalemir/hop/Onboardingv2/internal/service/domain"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	AuthorService domain.AuthorServiceInterface
}

func NewHandler(authorService domain.AuthorServiceInterface) *Handler {
	return &Handler{AuthorService: authorService}
}

func (h *Handler) RegisterAuthor(w http.ResponseWriter, r *http.Request) {
	var author entities.Author

	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := h.AuthorService.RegisterAuthor(r.Context(), &author); err != nil {
		http.Error(w, "Failed to register author", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(author)
}
func (h *Handler) GetAuthorByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	author, err := h.AuthorService.GetAuthorByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Author not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(author)
}
