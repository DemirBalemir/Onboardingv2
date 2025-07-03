package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/demirbalemir/hop/Onboardingv2/internal/db"
	server "github.com/demirbalemir/hop/Onboardingv2/internal/server/http/handler"
	"github.com/demirbalemir/hop/Onboardingv2/internal/service/domain"
	"github.com/demirbalemir/hop/Onboardingv2/internal/storage/postgres"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	// Connect to DB
	dbPool := db.NewPostgresConnection(dsn)
	defer dbPool.Close()

	// Initialize Repositories
	repo := postgres.NewRepository(dbPool)

	// Initialize Services
	authorService := domain.NewAuthorService(repo.Author, repo.Book)
	bookService := domain.NewBookService(repo.Book)

	// Start HTTP Server
	server.StartServer(authorService, bookService)
}
