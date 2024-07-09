package router

import (
	"fmt"
	"net/http"

	"goFinancialChat/features/chat"
	// Import other feature packages as needed
)

func init() {
	fmt.Println("Initializing router package...")
}

// SetupMux creates and configures the main router for the application
func SetupMux() *http.ServeMux {
	fmt.Println("Setting up router...")
	mux := http.NewServeMux()

	// Register routes
	registerStaticRoutes(mux)
	registerAPIRoutes(mux)
	registerWebRoutes(mux)

	fmt.Println("Router setup complete.")
	return mux
}

// registerStaticRoutes sets up routes for serving static files
func registerStaticRoutes(mux *http.ServeMux) {
	// Serve static files from the "static" directory
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	fmt.Println("Registered static file routes.")
}

// registerAPIRoutes sets up routes for API endpoints
func registerAPIRoutes(mux *http.ServeMux) {
	// Chat API route (now using the new handler)
	mux.HandleFunc("POST /api/chat", chat.NewHandleChat)

	fmt.Println("Registered API routes.")
}

// registerWebRoutes sets up routes for web pages
func registerWebRoutes(mux *http.ServeMux) {
	// Serve the main page
	mux.HandleFunc("GET /", chat.ServeChatPage)

	fmt.Println("Registered web routes.")
}