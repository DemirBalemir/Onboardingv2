/*
this file is the interface of the storage layer
it contains all the method signatures that
book and author use

*/

package storage

import (
	"context"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
)

type BookRepository interface {
	FindAll(ctx context.Context) ([]*entities.Book, error)
	FindById(ctx context.Context, id int) (*entities.Book, error)
	Create(ctx context.Context, book *entities.Book) error
	Update(ctx context.Context, book *entities.Book) error
	Delete(ctx context.Context, id int) error
}

type AuthorRepository interface {
	FindByID(ctx context.Context, id int) (*entities.Author, error)
	Create(ctx context.Context, author *entities.Author) error
}

type Repository struct {
	Book   BookRepository
	Author AuthorRepository
}
