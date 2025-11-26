package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/ecetinerdem/forseerv2/internal/env"
)

func New(addr string, maxOpenConn int, maxIdleConn int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)

	if err != nil {
		return nil, err
	}

	duration := env.GetDuration("MAX_IDLE_TIME", maxIdleTime)

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxIdleTime(time.Duration(duration))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
