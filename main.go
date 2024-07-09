package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"goFinancialChat/features/chat"
)

func main() {
	http.Handle("/", templ.Handler(chat.IndexPage("Chat")))
	http.HandleFunc("/api/chat", handleChat)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	message := r.FormValue("message")
	// Here you would process the message and generate a response
	response := fmt.Sprintf("You said: %s", message)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}