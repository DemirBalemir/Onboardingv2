package domain

import (
	"context"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/stretchr/testify/mock"
)

type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) GetAllBooks(ctx context.Context) ([]*entities.Book, error) {
	args := m.Called(ctx)

	books, _ := args.Get(0).([]*entities.Book)
	err := args.Error(1)
	return books, err
}

func (m *MockBookService) GetBookByID(ctx context.Context, id int) (*entities.Book, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Book), args.Error(1)
}

func (m *MockBookService) AddBook(ctx context.Context, book *entities.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockBookService) UpdateBook(ctx context.Context, book *entities.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockBookService) RemoveBook(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookService) SearchGoogleBooks(ctx context.Context, title string) ([]entities.GoogleBook, error) {
	args := m.Called(ctx, title)

	books, _ := args.Get(0).([]entities.GoogleBook)
	err := args.Error(1)
	return books, err
}
