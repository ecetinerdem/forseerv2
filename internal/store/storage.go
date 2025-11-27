package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Storage struct {
	Users interface {
		Create(context.Context, *User) error
	}
	Portfolio interface {
		Create(context.Context, *sql.Tx, *Portfolio) error
		CreatePortfolioWithStocks(context.Context, *Portfolio) error
		GetPortfolios(context.Context, int64) ([]*Portfolio, error)
		GetPortfolioByID(context.Context, int64, int64) (*Portfolio, error)
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
