package store

import (
	"context"
	"database/sql"
)

type Portfolio struct {
	ID        int64   `json:"id"`
	UserId    int64   `json:"user_id"`
	Name      string  `json:"name"`
	Stocks    []Stock `json:"stocks"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type Stock struct {
	ID           int64   `json:"id"`
	PortfolioID  int64   `json:"portfolio_id"`
	Symbol       string  `json:"symbol"`
	Shares       float64 `json:"shares"`
	AveragePrice float64 `json:"average_price"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type PortfolioStore struct {
	db *sql.DB
}

func (ps *PortfolioStore) Create(ctx context.Context, tx *sql.Tx, portfolio *Portfolio) error {
	query := `
		INSERT INTO portfolios (user_id, name, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`

	err := tx.QueryRowContext(ctx, query, portfolio.UserId, portfolio.Name).Scan(
		&portfolio.ID,
	)

	if err != nil {
		return err
	}

	for i := range portfolio.Stocks {
		stock := &portfolio.Stocks[i]
		stock.PortfolioID = portfolio.ID

		stockQuery := `
			INSERT INTO portfolio_stocks (portfolio_id, symbol, shares, average_price, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
		`
		_, err := tx.ExecContext(ctx, stockQuery, stock.PortfolioID, stock.Symbol, stock.Shares, stock.AveragePrice)

		if err != nil {
			return err
		}
	}

	return nil
}
