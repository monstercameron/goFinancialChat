package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"goFinancialChat/utils"
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

	// Sanity check 3: Create required tables
	if err := createTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Create a test user
	if err := createTestUser(); err != nil {
		log.Fatalf("Failed to create test user: %v", err)
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

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE,
			username TEXT,
			passphrase TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS conversations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			user_message TEXT,
			ai_response TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	return nil
}

func createTestUser() error {
	user, err := utils.CreateUser("test@gmail.com", "TestUser", "verylongpassword")
	if err != nil {
		return fmt.Errorf("failed to create test user: %v", err)
	}

	_, err = db.Exec("INSERT OR IGNORE INTO users (email, username, passphrase) VALUES (?, ?, ?)",
		user.Email, user.Username, user.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to insert test user: %v", err)
	}

	fmt.Println("Test user created or already exists.")
	return nil
}

func GetUserByEmail(email string) (*utils.User, error) {
	var user utils.User
	err := db.QueryRow("SELECT id, email, username, passphrase FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.Username, &user.Passphrase)
	if err != nil {
		return nil, err
	}
	return &user, nil
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
func SaveConversation(userID int, userMessage, aiResponse string) error {
	_, err := db.Exec("INSERT INTO conversations (user_id, user_message, ai_response) VALUES (?, ?, ?)",
		userID, userMessage, aiResponse)
	if err != nil {
		return fmt.Errorf("error saving conversation: %v", err)
	}
	return nil
}

// GetConversations retrieves all conversations from the database
func GetConversations(userID int) ([]Conversation, error) {
	rows, err := db.Query("SELECT id, user_message, ai_response, timestamp FROM conversations WHERE user_id = ? ORDER BY timestamp DESC", userID)
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

type Conversation struct {
	ID           int
	UserMessage  string
	AIResponse   string
	Timestamp    string
}
