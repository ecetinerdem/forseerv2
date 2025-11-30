package store

import (
	"database/sql"
)

type Stock struct {
	ID           int64   `json:"id"`
	PortfolioID  int64   `json:"portfolio_id"`
	Symbol       string  `json:"symbol"`
	Shares       float64 `json:"shares"`
	AveragePrice float64 `json:"average_price"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type StockStore struct {
	db *sql.DB
}
