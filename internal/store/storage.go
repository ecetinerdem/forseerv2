package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound        = errors.New("resource not found")
	ErrVersionConflict = errors.New("resource version conflict")
	QueryTimeOut       = time.Second * 5
)

type Storage struct {
	Users interface {
		Create(context.Context, *User) error
		CreateAndInvite(context.Context, *User, string) error
		GetUserByID(context.Context, int64) (*User, error)
	}
	Portfolio interface {
		Create(context.Context, *sql.Tx, *Portfolio) error
		CreatePortfolioWithStocks(context.Context, *Portfolio) error
		GetPortfolios(context.Context, int64) ([]*Portfolio, error)
		SearchPortfoliosByName(context.Context, int64, string) ([]*Portfolio, error)
		GetPortfolioByID(context.Context, int64) (*Portfolio, error)
		UpdatePortfolio(context.Context, *Portfolio) (*Portfolio, error)
		DeletePortfolio(context.Context, int64) error
	}
	Stock interface {
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:     &UserStore{db},
		Portfolio: &PortfolioStore{db},
	}
}

func withTX(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
