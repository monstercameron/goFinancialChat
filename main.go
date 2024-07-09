package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goFinancialChat/features/router"

	"github.com/joho/godotenv"
)

func init() {
	fmt.Println("Initializing main package...")

	// Load the .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found. Using default values or environment variables.")
	} else {
		fmt.Println("Successfully loaded .env file.")
	}
}

func main() {
	fmt.Println("Starting Financial Chat application...")

	// Set up routing
	fmt.Println("Setting up routes...")
	routes := router.SetupMux()
	fmt.Println("Routes successfully set up.")

	// Read PORT from environment variable, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("No PORT specified in environment. Defaulting to", port)
	} else {
		fmt.Println("Using PORT from environment:", port)
	}

	// Create a new HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: routes,
	}

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Server is starting on http://localhost:%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Fatal error while serving HTTP: %v\n", err)
		}
	}()
	fmt.Println("Server is now running. Press Ctrl+C to shut down.")

	// Set up channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	sig := <-quit
	fmt.Printf("Received shutdown signal: %v\n", sig)

	fmt.Println("Initiating graceful shutdown...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
		log.Fatal("Server forced to shutdown")
	}

	fmt.Println("Server has shut down gracefully")
	fmt.Println("Financial Chat application has terminated")
}