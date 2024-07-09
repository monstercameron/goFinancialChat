package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "./database.db"
)

var db *sql.DB

func init() {
	var err error

	// Sanity check 1: Check if the database file exists
	if !databaseFileExists() {
		fmt.Println("Database file does not exist. Creating a new one.")
		if err := createDatabaseFile(); err != nil {
			log.Fatalf("Failed to create database file: %v", err)
		}
	}

	// Open the database
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Sanity check 2: Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Successfully connected to the database.")

	// Sanity check 3: Check if the required table exists
	if !tableExists("conversations") {
		fmt.Println("Conversations table does not exist. Creating it.")
		if err := createConversationsTable(); err != nil {
			log.Fatalf("Failed to create conversations table: %v", err)
		}
	}

	fmt.Println("Database initialization completed successfully.")
}

func databaseFileExists() bool {
	_, err := os.Stat(dbPath)
	return !os.IsNotExist(err)
}

func createDatabaseFile() error {
	file, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database file: %v", err)
	}
	file.Close()
	return nil
}

func tableExists(tableName string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	var name string
	err := db.QueryRow(query, tableName).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Printf("Error checking if table exists: %v", err)
		return false
	}
	return true
}

func createConversationsTable() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS conversations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_message TEXT,
			ai_response TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating conversations table: %v", err)
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// SaveConversation saves a conversation to the database
func SaveConversation(userMessage, aiResponse string) error {
	_, err := db.Exec("INSERT INTO conversations (user_message, ai_response) VALUES (?, ?)", userMessage, aiResponse)
	if err != nil {
		return fmt.Errorf("error saving conversation: %v", err)
	}
	return nil
}

// GetConversations retrieves all conversations from the database
func GetConversations() ([]Conversation, error) {
	rows, err := db.Query("SELECT id, user_message, ai_response, timestamp FROM conversations ORDER BY timestamp DESC")
	if err != nil {
		return nil, fmt.Errorf("error querying conversations: %v", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var c Conversation
		err := rows.Scan(&c.ID, &c.UserMessage, &c.AIResponse, &c.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("error scanning conversation row: %v", err)
		}
		conversations = append(conversations, c)
	}

	return conversations, nil
}

// Conversation represents a single conversation entry
type Conversation struct {
	ID           int
	UserMessage  string
	AIResponse   string
	Timestamp    string
}