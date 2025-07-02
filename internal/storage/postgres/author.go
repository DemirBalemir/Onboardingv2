package postgres

import (
	"context"
	"fmt"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/jackc/pgx/v5"
)

type Author struct {
	db PgxIface
}

func NewAuthorRepository(pool PgxIface) *Author {
	return &Author{db: pool}
}

func (a *Author) FindByID(ctx context.Context, id int) (*entities.Author, error) {

	query :=
		`
		SELECT
			id,
			name,
			bio,
			birthdate
		FROM
			authors
		WHERE
			id = $1 -- Use $1 for the first parameter in pgx
	`
	author := &entities.Author{}
	err := a.db.QueryRow(ctx, query, id).Scan(
		&author.ID,
		&author.Name,
		&author.Bio,
		&author.BirthDate,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// If no row was found, return nil for the author and a specific error
			return nil, fmt.Errorf("author with ID %d not found: %w", id, err)
		}
		// For any other database error
		return nil, fmt.Errorf("failed to find author by ID %d: %w", id, err)
	}
	return author, nil
}
func (a *Author) Create(ctx context.Context, author *entities.Author) error {
	query := `
		INSERT INTO authors (name, bio, birthdate)
		VALUES ($1, $2, $3)
		RETURNING id -- This returns the auto-generated ID
	`

	// Use QueryRow because we expect to return the generated ID
	// Pass the fields of the 'author' struct as parameters to the query
	err := a.db.QueryRow(ctx, query, author.Name, author.Bio, author.BirthDate).Scan(&author.ID)
	if err != nil {
		return fmt.Errorf("failed to create author: %w", err)
	}

	// If successful, the author.ID field of the passed-in struct
	// will now be populated with the new ID from the database.
	return nil // No error
}
