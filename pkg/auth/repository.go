package auth

type UserRepository interface {
	CreateUser(user *User) error
	FindByUsername(username string) (*User, error)
	FindByID(id string) (*User, error)
}

type TokenRepository interface {
	CreateToken(token *Token) error
	FindToken(token string) (*Token, error)
	DeleteToken(token string) error
	DeleteExpired() error
}
