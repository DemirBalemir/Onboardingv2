package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	authorHandler "github.com/demirbalemir/hop/Onboardingv2/internal/server/http/handler/author"
	bookHandler "github.com/demirbalemir/hop/Onboardingv2/internal/server/http/handler/book"
	"github.com/demirbalemir/hop/Onboardingv2/internal/service/domain"
)

func NewRouter(
	authorService *domain.AuthorService,
	bookService *domain.BookService,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Create handlers
	authorH := authorHandler.NewHandler(authorService)
	bookH := bookHandler.NewHandler(bookService)

	// Register routes
	r.Route("/authors", func(r chi.Router) {
		authorHandler.RegisterRoutes(r, authorH)
	})

	r.Route("/books", func(r chi.Router) {
		bookHandler.RegisterRoutes(r, bookH)
	})

	return r
}
