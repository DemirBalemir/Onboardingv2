package domain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	results, err := svc.SearchGoogleBooks(
		context.WithValue(context.Background(), baseURLKey, server.URL),
		"harry potter",
	)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.GreaterOrEqual(t, len(results), 1)
	assert.Equal(t, "abc123", results[0].ID)
	assert.Equal(t, "Harry Potter and the Sorcerer's Stone", results[0].VolumeInfo.Title)
	assert.Equal(t, "J.K. Rowling", results[0].VolumeInfo.Authors[0])
}

// repoMock uses testify's mock to verify interactions with the repository.
type repoMock struct {
	mock.Mock
}

func (m *repoMock) FindAll(ctx context.Context) ([]*entities.Book, error) {
	args := m.Called(ctx)
	books, _ := args.Get(0).([]*entities.Book)
	return books, args.Error(1)
}

func (m *repoMock) FindById(ctx context.Context, id int) (*entities.Book, error) {
	args := m.Called(ctx, id)
	book, _ := args.Get(0).(*entities.Book)
	return book, args.Error(1)
}

func (m *repoMock) Create(ctx context.Context, book *entities.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *repoMock) Update(ctx context.Context, book *entities.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *repoMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// helper to create service with mock repository
func newServiceWithMock(repo *repoMock) *BookService {
	return &BookService{repo: repo, client: &http.Client{}}
}

func TestBookService_GetAllBooks(t *testing.T) {
	ctx := context.Background()
	expected := []*entities.Book{{ID: 1, Title: "Test"}}

	repo := &repoMock{}
	repo.On("FindAll", ctx).Return(expected, nil).Once()

	svc := newServiceWithMock(repo)

	books, err := svc.GetAllBooks(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, books)
	repo.AssertExpectations(t)
}

func TestBookService_GetBookByID(t *testing.T) {
	ctx := context.Background()
	expected := &entities.Book{ID: 2, Title: "Another"}

	repo := &repoMock{}
	repo.On("FindById", ctx, 2).Return(expected, nil).Once()

	svc := newServiceWithMock(repo)

	book, err := svc.GetBookByID(ctx, 2)
	assert.NoError(t, err)
	assert.Equal(t, expected, book)
	repo.AssertExpectations(t)
}

func TestBookService_AddBook(t *testing.T) {
	ctx := context.Background()
	b := &entities.Book{Title: "Create"}

	repo := &repoMock{}
	repo.On("Create", ctx, b).Return(nil).Once()

	svc := newServiceWithMock(repo)

	err := svc.AddBook(ctx, b)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookService_UpdateBook(t *testing.T) {
	ctx := context.Background()
	b := &entities.Book{ID: 3, Title: "Updated"}

	repo := &repoMock{}
	repo.On("Update", ctx, b).Return(nil).Once()

	svc := newServiceWithMock(repo)

	err := svc.UpdateBook(ctx, b)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookService_RemoveBook(t *testing.T) {
	ctx := context.Background()
	repo := &repoMock{}
	repo.On("Delete", ctx, 4).Return(nil).Once()

	svc := newServiceWithMock(repo)

	err := svc.RemoveBook(ctx, 4)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
