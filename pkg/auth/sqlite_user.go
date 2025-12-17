package auth

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/rhydori/logs"
)

type SQLiteUserRepo struct {
	db *sql.DB
}

func NewSQLiteUserRepo(db *sql.DB) *SQLiteUserRepo {
	return &SQLiteUserRepo{
		db: db,
	}
}

func (r *SQLiteUserRepo) CreateUser(u *User) error {
	_, err := r.db.Exec(`
	INSERT INTO users (id, username, password_hash, created_at)
	VALUES(?, ?, ?, ?)`, u.ID, u.Username, u.PasswordHash, u.CreatedAt,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return errors.New("Username already exists.")
		}
		logs.Error("CreateUser failed for username:", u.Username, "Error:", err)
		return errors.New("Failed to create user due to system error.")
	}
	return nil
}

func (r *SQLiteUserRepo) FindByUsername(username string) (*User, error) {
	row := r.db.QueryRow(`
	SELECT id, username, password_hash, created_at
	FROM users 
	WHERE username = ?`, username)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *SQLiteUserRepo) FindByID(id string) (*User, error) {
	row := r.db.QueryRow(`
	SELECT id, username, password_hash, created_at
	FROM users 
	WHERE id = ?`, id)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
