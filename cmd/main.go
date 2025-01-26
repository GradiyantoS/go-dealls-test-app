package main

import (
	"log"
	"net/http"

	"github.com/GradiyantoS/go-dealls-test-app/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Falling back to environment variables.")
	}
	router := routes.SetupRouter()
	// Start the server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
