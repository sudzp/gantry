// Package main is the entry point for the Gantry CI/CD server
// Package main is the entry point for the Gantry CI/CD server
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
	defer func() { _ = srv.Cleanup() }()

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
	log.Println("")
	log.Println("   ╔════════════════════════════════════════════════════════════╗")
	log.Println("   ║                                                            ║")
	log.Println("   ║   ██████╗  █████╗ ███╗   ██╗████████╗██████╗ ██╗   ██╗   ║")
	log.Println("   ║  ██╔════╝ ██╔══██╗████╗  ██║╚══██╔══╝██╔══██╗╚██╗ ██╔╝   ║")
	log.Println("   ║  ██║  ███╗███████║██╔██╗ ██║   ██║   ██████╔╝ ╚████╔╝    ║")
	log.Println("   ║  ██║   ██║██╔══██║██║╚██╗██║   ██║   ██╔══██╗  ╚██╔╝     ║")
	log.Println("   ║  ╚██████╔╝██║  ██║██║ ╚████║   ██║   ██║  ██║   ██║      ║")
	log.Println("   ║   ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝   ╚═╝      ║")
	log.Println("   ║                                                            ║")
	log.Println("   ║          🚀 Lightweight CI/CD Platform v1.0                ║")
	log.Println("   ║                                                            ║")
	log.Println("   ╚════════════════════════════════════════════════════════════╝")
	log.Println("")
}
