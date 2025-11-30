package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	IsActive  bool     `json:"is_active"`
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

func (us *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (first_name, last_name, username, email, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, is_active, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	err := tx.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Username, user.Email, user.Password).Scan(
		&user.ID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (us *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationEXP time.Duration) error {

	return withTX(us.db, ctx, func(tx *sql.Tx) error {
		err := us.Create(ctx, tx, user)
		if err != nil {
			return err
		}

		err = us.createUserInvitation(ctx, tx, token, invitationEXP, user.ID)
		if err != nil {
			return err
		}
		return nil
	})

}

func (us *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationEXP time.Duration, userID int64) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expiry)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationEXP))

	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, username, email, password, is_active, created_at, updated_at
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
		&user.IsActive,
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
