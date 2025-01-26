package integration_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	cwd, _ := os.Getwd()
	log.Println("Current Working Directory:", cwd)
	rootDir, _ := filepath.Abs("../..") // Adjust based on your directory structure
	err := godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Println("Warning Integration Test: .env file not found, using environment variables instead.")
	}

	// Set fallback JWT_SECRET_KEY if not set
	if os.Getenv("JWT_SECRET_KEY") == "" {
		os.Setenv("JWT_SECRET_KEY", "fallback-secret-key")
	}

	log.Println("JWT_SECRET_KEY:", os.Getenv("JWT_SECRET_KEY"))

	// Run tests
	os.Exit(m.Run())
}
