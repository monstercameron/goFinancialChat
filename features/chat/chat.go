package chat

import (
	"fmt"
	"net/http"
)

func init() {
	fmt.Println("chat package init")
}


func HandleChat(w http.ResponseWriter, r *http.Request) {
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