package main

import (
	"fmt"
	"log"
	"net/http"
	"goFinancialChat/features/router"
)

func init() {
	fmt.Println("main package init")
}

func main() {
	fmt.Println("Starting server")
	mux := router.SetupMux()

	// start listening
	http.Handle("/", mux)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
