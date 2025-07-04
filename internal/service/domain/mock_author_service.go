package domain

import (
	"context"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/stretchr/testify/mock"
)

type MockAuthorService struct {
	mock.Mock
}

func (m *MockAuthorService) RegisterAuthor(ctx context.Context, author *entities.Author) error {
	args := m.Called(ctx, author)
	return args.Error(0)
}

func (m *MockAuthorService) GetAuthorByID(ctx context.Context, id int) (*entities.Author, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Author), args.Error(1)
}
