package domain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/stretchr/testify/assert"
)

// Dummy repo that does nothing (we're testing API logic only)
type mockBookRepo struct{}

func (m *mockBookRepo) FindAll(ctx context.Context) ([]*entities.Book, error)        { return nil, nil }
func (m *mockBookRepo) FindById(ctx context.Context, id int) (*entities.Book, error) { return nil, nil }
func (m *mockBookRepo) Create(ctx context.Context, book *entities.Book) error        { return nil }
func (m *mockBookRepo) Update(ctx context.Context, book *entities.Book) error        { return nil }
func (m *mockBookRepo) Delete(ctx context.Context, id int) error                     { return nil }

func TestSearchGoogleBooks(t *testing.T) {
	fakeResponse := `{
        "items": [{
            "id": "abc123",
            "volumeInfo": {
                "title": "Harry Potter and the Sorcerer's Stone",
                "authors": ["J.K. Rowling"],
                "description": "A young wizard begins his journey."
            }
        }]
    }`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fakeResponse))
	}))
	defer server.Close()

	svc := &BookService{
		repo:   &mockBookRepo{},
		client: server.Client(), // Inject mock HTTP client
	}

	results, err := svc.SearchGoogleBooks(context.Background(), "harry potter")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "abc123", results[0].ID)
	assert.Equal(t, "Harry Potter and the Sorcerer's Stone", results[0].VolumeInfo.Title)
	assert.Equal(t, "J.K. Rowling", results[0].VolumeInfo.Authors[0])
}
