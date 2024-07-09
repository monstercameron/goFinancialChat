package router

import (
	"fmt"
	"net/http"

	"goFinancialChat/features/chat"
)

func init() {
	fmt.Println("router package init")
}

func SetupMux() *http.ServeMux {
	fmt.Println("Setting up router")
	// create a new mux
	mux := http.NewServeMux()

	// register the chat handler
	mux.HandleFunc("/chat", chat.HandleChat)

	// return the mux
	return mux
}
