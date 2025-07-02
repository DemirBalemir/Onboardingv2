// internal/service/service.go
package service

import (
	"context"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
)

type BookService interface {
	GetAllBooks(ctx context.Context) ([]*entities.Book, error)
	GetBookByID(ctx context.Context, id int) (*entities.Book, error)
	AddBook(ctx context.Context, book *entities.Book) error
	UpdateBook(ctx context.Context, book *entities.Book) error
	RemoveBook(ctx context.Context, id int) error
	SearchGoogleBooks(ctx context.Context, title string) ([]entities.GoogleBook, error)
}

type AuthorService interface {
	GetAuthorByID(ctx context.Context, id int) (*entities.Author, error)
	RegisterAuthor(ctx context.Context, author *entities.Author) error
}
type Service struct {
	Book   BookService
	Author AuthorService
}
