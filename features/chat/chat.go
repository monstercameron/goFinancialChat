package chat

import (
	"fmt"
	"net/http"
	"goFinancialChat/utils"
	"goFinancialChat/database"
	"encoding/json"
)

func init() {
	fmt.Println("chat package init")
}

// ServeChatPage handles GET requests to serve the chat page
func ServeChatPage(w http.ResponseWriter, r *http.Request) {
	// Retrieve last 5 conversations from the database
	conversations, err := database.GetConversations()
	if err != nil {
		fmt.Printf("Error retrieving conversations: %v\n", err)
		http.Error(w, "Error loading chat history", http.StatusInternalServerError)
		return
	}

	// Limit to last 5 conversations
	if len(conversations) > 5 {
		conversations = conversations[:5]
	}

	// Render the chat page with history
	ChatPageWithHistory(conversations).Render(r.Context(), w)
}

// HandleChat handles POST requests for chat messages and uses the database
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

	// Save the conversation to the database
	err = database.SaveConversation(message, aiResponse)
	if err != nil {
		fmt.Printf("Error saving conversation: %v\n", err)
		// Note: We're not returning an error to the user here, but you might want to handle this differently
	}

	// Render chat bubbles
	ChatBubbles(message, aiResponse).Render(r.Context(), w)
}

// Ternary is a helper function to mimic the ternary operator
func Ternary(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}