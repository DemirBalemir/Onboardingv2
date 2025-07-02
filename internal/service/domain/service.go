package domain

import (
	"github.com/demirbalemir/hop/Onboardingv2/internal/service"
	"github.com/demirbalemir/hop/Onboardingv2/internal/storage"
)

func NewService(repositories *storage.Repository) *service.Service {
	return &service.Service{
		Book:   NewBookService(repositories.Book),
		Author: NewAuthorService(repositories.Author, repositories.Book),
	}
}
