package postgres

import (
	"context"
	"fmt"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Book struct {
	db PgxIface
}

func NewBookRepository(db PgxIface) *Book {
	return &Book{db: db}
}

type PgxIface interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

func (b *Book) FindAll(ctx context.Context) ([]*entities.Book, error) {

	query := `
	SELECT
		id,
		title,
		description,
		published_at,
		author_id,
		price
	FROM
		books
	ORDER BY published_at DESC -- Good to have an ORDER BY even without pagination
	`

	rows, err := b.db.Query(ctx, query) // Use 'b.db' now
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close() // ALWAYS defer rows.Close()

	books := make([]*entities.Book, 0)
	for rows.Next() {
		book := &entities.Book{}
		// Scan all columns in the order they are selected
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.PublishedAt,
			&book.AuthorID,
			&book.Price,
		); err != nil {
			return nil, fmt.Errorf("failed to scan book row: %w", err)
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return books, nil

}

func (b *Book) FindById(ctx context.Context, id int) (*entities.Book, error) {

	query :=
		`
		SELECT
			id,
			title,
			description,
			published_at,
			author_id,
			price
		FROM
			books
		WHERE
			id = $1
	`

	book := &entities.Book{}
	err := b.db.QueryRow(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Description,
		&book.PublishedAt,
		&book.AuthorID,
		&book.Price,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// If no row was found, return nil for the author and a specific error
			return nil, fmt.Errorf("author with ID %d not found: %w", id, err)
		}
		// For any other database error
		return nil, fmt.Errorf("failed to find author by ID %d: %w", id, err)
	}

	return book, nil
}

func (b *Book) Create(ctx context.Context, book *entities.Book) error {
	query := `
    INSERT INTO books (title, description, published_at, author_id, price)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id -- This returns the auto-generated ID
	`
	err := b.db.QueryRow(ctx, query, book.Title, book.Description, book.PublishedAt, book.AuthorID, book.Price).Scan(&book.ID)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}
func (b *Book) Delete(ctx context.Context, id int) error {
	query :=
		`
		DELETE 
		
		FROM
			books
		WHERE
			id = $1 -- Use $1 for the first parameter in pgx
	`

	cmdTag, err := b.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book with ID %d: %w", id, err) // Corrected error message
	}

	// Check if any row was actually deleted (i.e., if the book existed)
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("book with ID %d not found for delete", id)
	}

	return nil

}

// internal/storage/postgres/book.go

// Update modifies an existing book's details in the database.
// It uses book.ID to identify the record to update.
func (b *Book) Update(ctx context.Context, book *entities.Book) error {
	query := `
        UPDATE books
        SET
            title = $1,
            description = $2,
            published_at = $3,
            author_id = $4,
            price = $5
        WHERE
            id = $6 -- The ID of the book to update
    `

	// Use Exec for UPDATE operations, as it doesn't return rows of data
	cmdTag, err := b.db.Exec(
		ctx,
		query,
		book.Title,
		book.Description,
		book.PublishedAt,
		book.AuthorID,
		book.Price,
		book.ID, // This is the value for $6 in the WHERE clause
	)
	if err != nil {
		return fmt.Errorf("failed to update book with ID %d: %w", book.ID, err)
	}

	// Check if any row was actually updated (i.e., if the book existed)
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("book with ID %d not found for update", book.ID)
	}

	return nil
}
