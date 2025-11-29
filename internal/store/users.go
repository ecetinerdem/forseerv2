package store

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &plainPassword
	p.hash = hash

	return nil
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

func (us *UserStore) CreateAndInvite(ctx context.Context, user *User, token string) error {

	return nil
}

func (us *UserStore) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	user := &User{}

	err := us.db.QueryRowContext(ctx, query, userID).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
