// internal/storage/postgres/storage.go
// why do ı need this file
package postgres

import (
	"github.com/demirbalemir/hop/Onboardingv2/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRepository(db *pgxpool.Pool) *storage.Repository {
	return &storage.Repository{
		Book:   NewBookRepository(db),   // postgres.Book implements storage.BookRepository
		Author: NewAuthorRepository(db), // postgres.Author implements storage.AuthorRepository
	}
}
