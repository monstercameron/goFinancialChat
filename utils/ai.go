package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

var (
	client *openai.Client
	model  string
)

func init() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	// Get the OpenAI API key from environment variables
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: OPENAI_API_KEY not set in environment")
	}

	// Get the OpenAI model from environment variables
	model = os.Getenv("OPENAI_MODEL")
	if model == "" {
		fmt.Println("Warning: OPENAI_MODEL not set in environment, defaulting to GPT-3.5-turbo")
		model = openai.GPT3Dot5Turbo
	}

	// Initialize the OpenAI client
	client = openai.NewClient(apiKey)
}

// GenerateAIResponse generates a response using the OpenAI API
func GenerateAIResponse(userMessage string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}

// IsAffirmativeResponse determines if a user's input is affirmative (yes) or negative (no)
func IsAffirmativeResponse(userInput string) (map[string]bool, error) {
	functionDefinition := openai.FunctionDefinition{
		Name: "determine_affirmative",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"isAffirmative": {
					"type": "boolean",
					"description": "True if the user's response is affirmative (yes), false if negative (no)"
				}
			},
			"required": ["isAffirmative"]
		}`),
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userInput,
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Determine if the user's response is affirmative (yes) or negative (no). Return only the JSON object with the result.",
				},
			},
			Functions: []openai.FunctionDefinition{functionDefinition},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("ChatCompletion error: %v", err)
	}

	if resp.Choices[0].Message.FunctionCall == nil {
		return nil, fmt.Errorf("No function call in the response")
	}

	var result map[string]bool
	err = json.Unmarshal([]byte(resp.Choices[0].Message.FunctionCall.Arguments), &result)
	if err != nil {
		return nil, fmt.Errorf("Error parsing function call arguments: %v", err)
	}

	return result, nil
}