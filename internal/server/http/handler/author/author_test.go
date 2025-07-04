package author

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	"github.com/demirbalemir/hop/Onboardingv2/internal/service/domain"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterAuthor(t *testing.T) {
	mockService := new(domain.MockAuthorService)
	handler := NewHandler(mockService)

	author := entities.Author{ID: 1, Name: "Test Author"}

	// Set expectation
	mockService.On("RegisterAuthor", mock.Anything, &author).Return(nil)

	body, _ := json.Marshal(author)
	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	// Set up router
	r := chi.NewRouter()
	r.Post("/authors", handler.RegisterAuthor)

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var result entities.Author
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "Test Author", result.Name)

	mockService.AssertExpectations(t)
}

func TestGetAuthorByID(t *testing.T) {
	mockService := new(domain.MockAuthorService)
	handler := NewHandler(mockService)

	author := &entities.Author{ID: 1, Name: "Test Author"}

	mockService.On("GetAuthorByID", mock.Anything, 1).Return(author, nil)

	req := httptest.NewRequest(http.MethodGet, "/authors/1", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/authors/{id}", handler.GetAuthorByID)

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result entities.Author
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Test Author", result.Name)

	mockService.AssertExpectations(t)
}

func TestGetAuthorByID_NotFound(t *testing.T) {
	mockService := new(domain.MockAuthorService)
	handler := NewHandler(mockService)

	mockService.On("GetAuthorByID", mock.Anything, 42).Return((*entities.Author)(nil), errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/authors/42", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/authors/{id}", handler.GetAuthorByID)

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockService.AssertExpectations(t)
}
