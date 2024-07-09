package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type User struct {
	ID        int
	Email     string
	Username  string
	Passphrase string
}

// GenerateUserID generates a unique user ID
func GenerateUserID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// CreateUser creates a new user
func CreateUser(email, username, passphrase string) (*User, error) {
	// In a real application, you'd want to hash the passphrase
	// and perform additional validation
	return &User{
		Email:     email,
		Username:  username,
		Passphrase: passphrase,
	}, nil
}

// AuthenticateUser authenticates a user
func AuthenticateUser(email, passphrase string) (*User, error) {
	// In a real application, you'd query the database and verify the passphrase
	// This is a placeholder implementation
	if email == "test@gmail.com" && passphrase == "verylongpassword" {
		return &User{
			Email:     email,
			Username:  "TestUser",
			Passphrase: passphrase,
		}, nil
	}
	return nil, fmt.Errorf("invalid credentials")
}