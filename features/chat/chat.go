package chat

import (
	"fmt"
	"net/http"
	"goFinancialChat/utils"
	"encoding/json"
)

func init() {
	fmt.Println("chat package init")
}

// ServeChatPage handles GET requests to serve the chat page
func ServeChatPage(w http.ResponseWriter, r *http.Request) {
	ChatPage().Render(r.Context(), w)
}

// Original HandleChat (commented out for reference)
/*
// HandleChat handles POST requests for chat messages
func HandleChat(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	
	// Generate AI response
	aiResponse, err := utils.GenerateAIResponse(message)
	if err != nil {
		fmt.Printf("Error generating AI response: %v\n", err)
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	// Render chat bubbles
	ChatBubble(message, true).Render(r.Context(), w)
	ChatBubble(aiResponse, false).Render(r.Context(), w)
}
*/

// NewHandleChat handles POST requests for chat messages and tests IsAffirmativeResponse
func NewHandleChat(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	
	// Check if the message is affirmative
	affirmativeResponse, err := utils.IsAffirmativeResponse(message)
	if err != nil {
		fmt.Printf("Error checking affirmative response: %v\n", err)
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	// Prepare the AI response
	var aiResponse string
	if affirmativeResponse["isAffirmative"] {
		aiResponse = "I understand that your response is affirmative (yes)."
	} else {
		aiResponse = "I understand that your response is negative (no)."
	}

	// Add the raw JSON response for debugging
	jsonResponse, _ := json.MarshalIndent(affirmativeResponse, "", "  ")
	aiResponse += fmt.Sprintf("\n\nDebug Info:\nRaw JSON response:\n%s", string(jsonResponse))

	// Render chat bubbles
	ChatBubble(message, true).Render(r.Context(), w)
	ChatBubble(aiResponse, false).Render(r.Context(), w)
}

// Ternary is a helper function to mimic the ternary operator
func Ternary(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}