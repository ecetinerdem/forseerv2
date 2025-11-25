package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Users interface {
		Create(context.Context, *User) error
	}
	Portfolio interface {
		Create(context.Context, *sql.Tx) error
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
