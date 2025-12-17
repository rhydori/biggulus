package auth

import (
	"database/sql"
	"time"
)

type SQLiteTokenRepo struct {
	db *sql.DB
}

func NewSQLiteTokenRepo(db *sql.DB) *SQLiteTokenRepo {
	return &SQLiteTokenRepo{
		db: db,
	}
}

func (r *SQLiteTokenRepo) CreateToken(t *Token) error {
	_, err := r.db.Exec(`
	INSERT INTO tokens (token, user_id, expires_at, created_at)
	VALUES(?, ?, ?, ?)`, t.Value, t.UserID, t.ExpiresAt, t.CreatedAt,
	)
	return err
}

func (r *SQLiteTokenRepo) FindToken(token string) (*Token, error) {
	row := r.db.QueryRow(`
	SELECT token, user_id, expires_at, created_at
	FROM tokens
	WHERE token = ?`, token)

	var t Token
	err := row.Scan(&t.Value, &t.UserID, &t.ExpiresAt, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *SQLiteTokenRepo) DeleteToken(token string) error {
	_, err := r.db.Exec(`
	DELETE FROM tokens
	WHERE token = ?`, token)
	return err
}

func (r *SQLiteTokenRepo) DeleteExpired() error {
	_, err := r.db.Exec(`
	DELETE FROM tokens
	WHERE expires_at < ?`, time.Now())
	return err
}
