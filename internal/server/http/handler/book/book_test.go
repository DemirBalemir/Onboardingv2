package book_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/demirbalemir/hop/Onboardingv2/internal/entities"
	book "github.com/demirbalemir/hop/Onboardingv2/internal/server/http/handler/book"
	"github.com/demirbalemir/hop/Onboardingv2/internal/service/domain"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(h *book.Handler) *chi.Mux {
	r := chi.NewRouter()
	book.RegisterRoutes(r, h)
	return r
}

func TestBookHandlers(t *testing.T) {
	mockService := new(domain.MockBookService)
	h := book.NewHandler(mockService)
	r := setupRouter(h)

	mockBook := &entities.Book{ID: 1, Title: "Go 101"}
	mockBookList := []*entities.Book{mockBook}
	mockGoogleBooks := []entities.GoogleBook{{
		ID: "g1",
		VolumeInfo: struct {
			Title       string   `json:"title"`
			Authors     []string `json:"authors"`
			Description string   `json:"description"`
		}{
			Title:       "Go by Google",
			Authors:     []string{"Google Inc."},
			Description: "A great book by Google",
		},
	}}

	tests := []struct {
		name       string
		method     string
		url        string
		body       interface{}
		mockSetup  func()
		expectCode int
	}{
		{
			name:   "GetAllBooks - success",
			method: http.MethodGet,
			url:    "/",
			mockSetup: func() {
				mockService.On("GetAllBooks", mock.Anything).Return(mockBookList, nil).Once()
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "GetBookByID - found",
			method: http.MethodGet,
			url:    "/1",
			mockSetup: func() {
				mockService.On("GetBookByID", mock.Anything, 1).Return(mockBook, nil).Once()
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "AddBook - success",
			method: http.MethodPost,
			url:    "/",
			body:   mockBook,
			mockSetup: func() {
				mockService.On("AddBook", mock.Anything, mockBook).Return(nil).Once()
			},
			expectCode: http.StatusCreated,
		},
		{
			name:   "UpdateBook - success",
			method: http.MethodPut,
			url:    "/",
			body:   mockBook,
			mockSetup: func() {
				mockService.On("UpdateBook", mock.Anything, mockBook).Return(nil).Once()
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "DeleteBook - success",
			method: http.MethodDelete,
			url:    "/1",
			mockSetup: func() {
				mockService.On("RemoveBook", mock.Anything, 1).Return(nil).Once()
			},
			expectCode: http.StatusNoContent,
		},
		{
			name:   "SearchGoogleBooks - success",
			method: http.MethodGet,
			url:    "/search/google?title=go",
			mockSetup: func() {
				mockService.On("SearchGoogleBooks", mock.Anything, "go").Return(mockGoogleBooks, nil).Once()
			},
			expectCode: http.StatusOK,
		},
		{
			name:       "SearchGoogleBooks - missing title",
			method:     http.MethodGet,
			url:        "/search/google",
			mockSetup:  func() {},
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil // reset previous calls
			tc.mockSetup()

			var req *http.Request
			if tc.body != nil {
				bodyBytes, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(tc.method, tc.url, bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.url, nil)
			}

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectCode, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}
