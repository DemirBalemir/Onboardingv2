package domain

import (
	"context"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/demirbalemir/hop/Onboardingv2/internal/storage"
)

type AuthorService struct {
	repo     storage.AuthorRepository
	bookRepo storage.BookRepository
}
type AuthorServiceInterface interface {
	RegisterAuthor(ctx context.Context, author *entities.Author) error
	GetAuthorByID(ctx context.Context, id int) (*entities.Author, error)
}

var _ AuthorServiceInterface = (*AuthorService)(nil)

func NewAuthorService(authorRepo storage.AuthorRepository, bookRepo storage.BookRepository) *AuthorService {
	return &AuthorService{
		repo:     authorRepo,
		bookRepo: bookRepo,
	}
}

func (s *AuthorService) GetAuthorByID(ctx context.Context, id int) (*entities.Author, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *AuthorService) RegisterAuthor(ctx context.Context, author *entities.Author) error {
	return s.repo.Create(ctx, author)
}
