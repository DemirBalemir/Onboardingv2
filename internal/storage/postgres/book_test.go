package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/demirbalemir/hop/Onboardingv2/internal/storage/postgres"
)

func setupMockRepo(t *testing.T) (pgxmock.PgxPoolIface, *postgres.Book, func()) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repo := postgres.NewBookRepository(mockPool)

	cleanup := func() {
		if err := mockPool.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		mockPool.Close()
	}
	return mockPool, repo, cleanup
}

func TestBookRepository_FindById(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		bookID       int
		mockSetup    func(mockPool pgxmock.PgxPoolIface)
		expectedBook *entities.Book
		expectedErr  error
	}{
		{
			name:   "should return a book when found",
			bookID: 1,
			mockSetup: func(mockPool pgxmock.PgxPoolIface) {
				expectedBook := &entities.Book{
					ID:          1,
					Title:       "Mocked Book 1",
					Description: "A description.",
					PublishedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					AuthorID:    101,
					Price:       19.99,
				}
				rows := pgxmock.NewRows([]string{"id", "title", "description", "published_at", "author_id", "price"}).
					AddRow(expectedBook.ID, expectedBook.Title, expectedBook.Description, expectedBook.PublishedAt, expectedBook.AuthorID, expectedBook.Price)

				mockPool.ExpectQuery(`SELECT id, title, description, published_at, author_id, price FROM books WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedBook: &entities.Book{
				ID:          1,
				Title:       "Mocked Book 1",
				Description: "A description.",
				PublishedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				AuthorID:    101,
				Price:       19.99,
			},
			expectedErr: nil,
		},
		{
			name:   "should return ErrNoRows when book not found",
			bookID: 999,
			mockSetup: func(mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectQuery(`SELECT id, title, description, published_at, author_id, price FROM books WHERE id = \$1`).
					WithArgs(999).
					WillReturnError(pgx.ErrNoRows) // Simulate no rows found
			},
			expectedBook: nil,
			expectedErr:  pgx.ErrNoRows,
		},
		{
			name:   "should return error for database query failure",
			bookID: 2,
			mockSetup: func(mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectQuery(`SELECT id, title, description, published_at, author_id, price FROM books WHERE id = \$1`).
					WithArgs(2).
					WillReturnError(errors.New("db connection lost")) // Simulate a generic DB error
			},
			expectedBook: nil,
			expectedErr:  errors.New("db connection lost"), // The expected error returned by your repo
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockPool, repo, cleanup := setupMockRepo(t)
			defer cleanup()

			tc.mockSetup(mockPool)

			foundBook, err := repo.FindById(ctx, tc.bookID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				if errors.Is(tc.expectedErr, pgx.ErrNoRows) {
					assert.True(t, errors.Is(err, pgx.ErrNoRows), "expected error to be pgx.ErrNoRows")
				} else {
					assert.Contains(t, err.Error(), tc.expectedErr.Error())
				}
				assert.Nil(t, foundBook)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, foundBook)
				assert.Equal(t, tc.expectedBook.ID, foundBook.ID)
				assert.Equal(t, tc.expectedBook.Title, foundBook.Title)
				assert.Equal(t, tc.expectedBook.Description, foundBook.Description)
				assert.WithinDuration(t, tc.expectedBook.PublishedAt, foundBook.PublishedAt, time.Second)
				assert.Equal(t, tc.expectedBook.AuthorID, foundBook.AuthorID)
				assert.Equal(t, tc.expectedBook.Price, foundBook.Price)
			}
		})
	}
}

func TestBookRepository_Create(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		bookToCreate  *entities.Book
		mockSetup     func(mockPool pgxmock.PgxPoolIface, book *entities.Book)
		expectedError error
	}{
		{
			name: "should successfully create a book",
			bookToCreate: &entities.Book{
				Title:       "New Mocked Book",
				Description: "A freshly created book.",
				PublishedAt: time.Date(2024, 5, 15, 10, 0, 0, 0, time.UTC),
				AuthorID:    201,
				Price:       29.99,
			},
			mockSetup: func(mockPool pgxmock.PgxPoolIface, book *entities.Book) {
				// Expect an INSERT query returning the ID
				// Adjust regex to exactly match your INSERT query in postgres/book.go
				mockPool.ExpectQuery(`INSERT INTO books \(title, description, published_at, author_id, price\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
					WithArgs(book.Title, book.Description, book.PublishedAt, book.AuthorID, book.Price).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(5)) // Simulate returning new ID 5
			},
			expectedError: nil,
		},
		{
			name: "should return error on database insert failure",
			bookToCreate: &entities.Book{
				Title:       "Failing Book",
				Description: "Will fail.",
				PublishedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				AuthorID:    202,
				Price:       15.00,
			},
			mockSetup: func(mockPool pgxmock.PgxPoolIface, book *entities.Book) {
				mockPool.ExpectQuery(`INSERT INTO books \(title, description, published_at, author_id, price\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
					WithArgs(book.Title, book.Description, book.PublishedAt, book.AuthorID, book.Price).
					WillReturnError(errors.New("duplicate key error")) // Simulate a DB error
			},
			expectedError: errors.New("duplicate key error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockPool, repo, cleanup := setupMockRepo(t)
			defer cleanup()

			tc.mockSetup(mockPool, tc.bookToCreate)

			// Act
			err := repo.Create(ctx, tc.bookToCreate)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				assert.Equal(t, 0, tc.bookToCreate.ID, "ID should not be set on error")
			} else {
				assert.NoError(t, err)
				assert.True(t, tc.bookToCreate.ID > 0, "Expected ID to be set after successful creation")
				// Additional check: You could retrieve from an in-memory map here if you passed it around,
				// but for pure pgxmock, verifying ID is enough.
			}
		})
	}
}

func TestBookRepository_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		bookIDToDelete int
		mockSetup      func(mockPool pgxmock.PgxPoolIface)
		expectedError  error
	}{
		{
			name:           "should successfully delete a book",
			bookIDToDelete: 1,
			mockSetup: func(mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectExec(`DELETE FROM books WHERE id = \$1`).
					WithArgs(1).
					WillReturnResult(pgxmock.NewResult("DELETE", 1)) // Simulate 1 row affected
			},
			expectedError: nil,
		},
		{
			name:           "should return error if book not found for delete",
			bookIDToDelete: 999,
			mockSetup: func(mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectExec(`DELETE FROM books WHERE id = \$1`).
					WithArgs(999).
					WillReturnResult(pgxmock.NewResult("DELETE", 0)) // Simulate 0 rows affected
			},
			// Your repository's Delete method should return an error like "not found"
			// if RowsAffected is 0. Adjust this expected error string accordingly.
			expectedError: errors.New("book with ID 999 not found for delete"),
		},
		{
			name:           "should return error on database delete failure",
			bookIDToDelete: 2,
			mockSetup: func(mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectExec(`DELETE FROM books WHERE id = \$1`).
					WithArgs(2).
					WillReturnError(errors.New("db error during delete")) // Simulate a generic DB error
			},
			expectedError: errors.New("db error during delete"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockPool, repo, cleanup := setupMockRepo(t)
			defer cleanup()

			tc.mockSetup(mockPool)

			// Act
			err := repo.Delete(ctx, tc.bookIDToDelete)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TODO: Add similar tests for Update and FindAll methods
