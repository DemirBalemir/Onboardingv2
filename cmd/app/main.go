package main

import (
	"log"
	"os"

	"github.com/demirbalemir/hop/Onboardingv2/internal/db"
	"github.com/joho/godotenv"
	//"github.com/demirbalemir/hop/Onboardingv2/internal/storage/postgres"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
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
	//repo := postgres.NewRepository(dbPool)

	log.Println("Application setup complete")

	// Example placeholder
	// server := server.New(repo)
	// server.Run()
}
