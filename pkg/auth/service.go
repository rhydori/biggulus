package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("User not found.")
	ErrIncorrectPassword = errors.New("Incorrect password.")
	ErrInvalidToken      = errors.New("Invalid or Expired Token.")
	ErrHashPassword      = errors.New("Error hashing password.")
	ErrInvalidPayload    = errors.New("Invalid message payload.")
)

type AuthService struct {
	users  UserRepository
	tokens TokenRepository

	tokenTTL time.Duration
}

func NewService(u UserRepository, t TokenRepository) *AuthService {
	return &AuthService{
		users:    u,
		tokens:   t,
		tokenTTL: 24 * time.Hour,
	}
}

func (as *AuthService) newToken(userID string) (*Token, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	t := &Token{
		Value:     hex.EncodeToString(b),
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(as.tokenTTL),
	}
	return t, as.tokens.CreateToken(t)
}

func (as *AuthService) Register(msg []string) error {
	if len(msg) < 2 {
		return ErrInvalidPayload
	}
	username := msg[0]
	password := msg[1]
	if err := validateUser(username, password); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrHashPassword
	}
	u := &User{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}
	return as.users.CreateUser(u)
}

func (as *AuthService) Login(msg []string) (*Token, error) {
	if len(msg) < 2 {
		return nil, ErrInvalidPayload
	}
	username := msg[0]
	password := msg[1]
	u, err := as.users.FindByUsername(username)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, ErrIncorrectPassword
	}
	return as.newToken(u.ID)
}

func (as *AuthService) Logout(msg []string) error {
	if len(msg) < 1 {
		return ErrInvalidPayload
	}
	token := msg[0]
	return as.tokens.DeleteToken(token)
}

var usernameRe = regexp.MustCompile(`^[A-Za-z0-9]+$`)

func validateUser(username, password string) error {
	switch {
	case username == "":
		return errors.New("Username cannot be empty.")
	case password == "":
		return errors.New("Password cannot be empty.")
	case len(username) < 4:
		return errors.New("Username must be at least 4 characters.")
	case len(password) < 8:
		return errors.New("Password must be at least 8 characters.")
	case len(username) > 16:
		return errors.New("Username must be less than 16 characters.")
	case len(password) > 20:
		return errors.New("Password must be less than 20 characters.")
	case !usernameRe.MatchString(username):
		return errors.New("Username contains invalid characters.")
	}
	return nil
}
