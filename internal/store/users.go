package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserStore struct {
	db *sql.DB
}

func (us *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (first_name, last_name, username, email, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := us.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Username, user.Email, user.Password).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
