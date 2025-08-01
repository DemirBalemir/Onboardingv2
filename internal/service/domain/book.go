package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/demirbalemir/hop/Onboardingv2/internal/storage"
)

type BookService struct {
	repo   storage.BookRepository
	client *http.Client
}

type ctxKey string

const baseURLKey ctxKey = "baseURL"

type GoogleBooksSearchResponse struct {
	Items []entities.GoogleBook `json:"items"`
}

type BookServiceInterface interface {
	GetAllBooks(ctx context.Context) ([]*entities.Book, error)
	GetBookByID(ctx context.Context, id int) (*entities.Book, error)
	AddBook(ctx context.Context, book *entities.Book) error
	UpdateBook(ctx context.Context, book *entities.Book) error
	RemoveBook(ctx context.Context, id int) error
	SearchGoogleBooks(ctx context.Context, title string) ([]entities.GoogleBook, error)
}

func NewBookService(repo storage.BookRepository) *BookService {
	return &BookService{
		repo:   repo,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]*entities.Book, error) {
	return s.repo.FindAll(ctx)
}

func (s *BookService) GetBookByID(ctx context.Context, id int) (*entities.Book, error) {
	return s.repo.FindById(ctx, id)
}

func (s *BookService) AddBook(ctx context.Context, book *entities.Book) error {
	return s.repo.Create(ctx, book)
}

func (s *BookService) UpdateBook(ctx context.Context, book *entities.Book) error {
	return s.repo.Update(ctx, book)
}

func (s *BookService) RemoveBook(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// ✅ Google Books API logic
//
// SearchGoogleBooks fetches a list of books from Google Books API by title
func (s *BookService) SearchGoogleBooks(ctx context.Context, title string) ([]entities.GoogleBook, error) {
	baseURL := "https://www.googleapis.com/books/v1/volumes"

	// Allow test to override base URL via context
	if custom := ctx.Value(baseURLKey); custom != nil {
		if str, ok := custom.(string); ok {
			baseURL = str
		}
	}

	url := fmt.Sprintf("%s?q=%s", baseURL, url.QueryEscape(title))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("google books API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result GoogleBooksSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Google Books API response: %w", err)
	}

	return result.Items, nil
}
