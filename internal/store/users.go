package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/lib/pq"
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
func (p *password) Compare(plainPassword string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plainPassword))
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
	err := tx.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Username, user.Email, user.Password.hash).Scan(
		&user.ID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Constraint {
			case "users_email_key":
				return ErrDuplicateEmail
			case "users_username_key":
				return ErrDuplicateUsername
			}
		}
		return err
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

func (us *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN user_invitations ui ON  u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	user := &User{}

	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
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

func (us *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, is_active = $3
		WHERE id = $4
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {

	query := `
		DELETE FROM user_invitations
		WHERE user_id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)

	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) Activate(ctx context.Context, token string) error {

	return withTX(us.db, ctx, func(tx *sql.Tx) error {
		user, err := us.getUserFromInvitation(ctx, tx, token)

		if err != nil {
			return err
		}

		user.IsActive = true

		err = us.update(ctx, tx, user)
		if err != nil {
			return err
		}

		err = us.deleteUserInvitations(ctx, tx, user.ID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (us *UserStore) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, username, email, password, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	user := &User{}

	err := us.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password.hash,
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

func (us *UserStore) DeleteUser(ctx context.Context, userID int64) error {
	return withTX(us.db, ctx, func(tx *sql.Tx) error {
		err := us.deleteUser(ctx, tx, userID)
		if err != nil {
			return err
		}

		err = us.deleteUserInvitations(ctx, tx, userID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (us *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, username, email, password, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	user := &User{}

	err := us.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password.hash,
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

func (us *UserStore) deleteUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)

	if err != nil {
		return err
	}

	return nil
}
