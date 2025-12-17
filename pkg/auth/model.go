package auth

import "time"

type User struct {
	ID           string
	Username     string
	Password     string
	PasswordHash string
	CreatedAt    time.Time
}

type Token struct {
	Value     string
	UserID    string
	ExpiresAt time.Time
	CreatedAt time.Time
}
