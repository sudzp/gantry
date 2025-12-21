package main

import (
	"log"
	"net/http"
	"os"

	"gantry/internal/api"
	"gantry/internal/server"
)

func main() {
	// Print banner
	printBanner()

	// Create server from environment variables
	srv, err := server.NewServerFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Cleanup()

	// Setup API handlers
	handler := api.NewHandler(srv)
	router := api.SetupRoutes(handler)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func printBanner() {
	log.Println("========================================")
	log.Println("    ____            __            ")
	log.Println("   / __ \\____ _____/ /________  __")
	log.Println("  / / / / __ `/ __  / ___/ / / /")
	log.Println(" / /_/ / /_/ / /_/ / /  / /_/ / ")
	log.Println(" \\____/\\__,_/\\__,_/_/   \\__, /  ")
	log.Println("                       /____/   ")
	log.Println("")
	log.Println("  Gantry CI/CD Platform")
	log.Println("========================================")
}
