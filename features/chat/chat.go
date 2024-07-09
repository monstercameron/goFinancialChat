package chat

import (
	"encoding/json"
	"fmt"
	"goFinancialChat/database"
	"goFinancialChat/utils"
	"net/http"
)

func init() {
	fmt.Println("chat package init")
}

// ServeChatPage handles GET requests to serve the chat page
func ServeChatPage(w http.ResponseWriter, r *http.Request) {
	// In a real application, you'd get the user from the session
	// For now, we'll use the test user
	user, err := database.GetUserByEmail("test@gmail.com")
	if err != nil {
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	conversations, err := database.GetConversations(user.ID)
	if err != nil {
		fmt.Printf("Error retrieving conversations: %v\n", err)
		http.Error(w, "Error loading chat history", http.StatusInternalServerError)
		return
	}

	if len(conversations) > 5 {
		conversations = conversations[:5]
	}

	ChatPageWithHistory(user, conversations).Render(r.Context(), w)
}

// HandleChat handles POST requests for chat messages and uses the database
func NewHandleChat(w http.ResponseWriter, r *http.Request) {
	// In a real application, you'd get the user from the session
	// For now, we'll use the test user
	user, err := database.GetUserByEmail("test@gmail.com")
	if err != nil {
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	message := r.FormValue("message")

	affirmativeResponse, err := utils.IsAffirmativeResponse(message)
	if err != nil {
		fmt.Printf("Error checking affirmative response: %v\n", err)
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	var aiResponse string
	if affirmativeResponse["isAffirmative"] {
		aiResponse = "I understand that your response is affirmative (yes)."
	} else {
		aiResponse = "I understand that your response is negative (no)."
	}

	jsonResponse, _ := json.MarshalIndent(affirmativeResponse, "", "  ")
	aiResponse += fmt.Sprintf("\n\nDebug Info:\nRaw JSON response:\n%s", string(jsonResponse))

	err = database.SaveConversation(user.ID, message, aiResponse)
	if err != nil {
		fmt.Printf("Error saving conversation: %v\n", err)
	}

	ChatBubbles(message, aiResponse).Render(r.Context(), w)
}

// Ternary is a helper function to mimic the ternary operator
func Ternary(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}
