package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/demirbalemir/hop/Onboardingv2/internal/service/domain"
)

func StartServer(authorService *domain.AuthorService, bookService *domain.BookService) {
	// Set up router
	router := NewRouter(authorService, bookService)

	// Start server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	fmt.Println("🚀 Server is running on http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
