package store

import (
	"context"
	"database/sql"
	"errors"
)

type Portfolio struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	Name      string  `json:"name"`
	Stocks    []Stock `json:"stocks"`
	Version   int     `json:"version"`
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

	err := tx.QueryRowContext(ctx, query, portfolio.UserID, portfolio.Name).Scan(
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

func (ps *PortfolioStore) CreatePortfolioWithStocks(ctx context.Context, portfolio *Portfolio) error {

	return withTX(ps.db, ctx, func(tx *sql.Tx) error {
		err := ps.Create(ctx, tx, portfolio)

		if err != nil {
			return err
		}
		return nil
	})
}

func (ps *PortfolioStore) GetPortfolios(ctx context.Context, userID int64) ([]*Portfolio, error) {

	query := `
		SELECT id, user_id, name, created_at, updated_at 
		FROM portfolios
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := ps.db.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var portfolios []*Portfolio

	for rows.Next() {
		var p Portfolio

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Name,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		portfolios = append(portfolios, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return portfolios, nil
}

func (ps *PortfolioStore) SearchPortfoliosByName(ctx context.Context, userId int64, searchParam string) ([]*Portfolio, error) {

	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM portfolios
		WHERE user_id = $1 AND name ILIKE $2
		ORDER BY updated_at DESC
	`

	searchPattern := "%" + searchParam + "%"

	rows, err := ps.db.QueryContext(ctx, query, userId, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	portfolios := []*Portfolio{}

	for rows.Next() {
		var p Portfolio

		err = rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Name,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		stockQuery := `
			SELECT id, portfolio_id, symbol, shares, average_price, created_at, updated_at
			FROM portfolio_stocks
			WHERE portfolio_id = $1
			ORDER BY symbol ASC
		`

		stockRows, err := ps.db.QueryContext(ctx, stockQuery, p.ID)

		if err != nil {
			return nil, err
		}

		p.Stocks = []Stock{}

		for stockRows.Next() {
			var stock Stock

			err = stockRows.Scan(
				&stock.ID,
				&stock.PortfolioID,
				&stock.Symbol,
				&stock.Shares,
				&stock.AveragePrice,
				&stock.CreatedAt,
				&stock.UpdatedAt,
			)

			if err != nil {
				stockRows.Close()
				return nil, err
			}
			stockRows.Close()
			p.Stocks = append(p.Stocks, stock)
		}
		if err := stockRows.Err(); err != nil {
			return nil, err
		}
		portfolios = append(portfolios, &p)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return portfolios, nil
}

func (ps *PortfolioStore) GetPortfolioByID(ctx context.Context, portfolioID int64, userID int64) (*Portfolio, error) {

	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM portfolios
		WHERE id=$1 AND user_id = $2
	`
	var p Portfolio

	err := ps.db.QueryRowContext(ctx, query, portfolioID, userID).Scan(
		&p.ID,
		&p.UserID,
		&p.Name,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	//Get stocks due to transaction for creation

	stockQuery := `
		SELECT id, portfolio_id, symbol, shares, average_price, created_at, updated_at
		FROM portfolio_stocks
		WHERE portfolio_id = $1
		ORDER BY symbol ASC
	`
	rows, err := ps.db.QueryContext(ctx, stockQuery, portfolioID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	p.Stocks = []Stock{}

	for rows.Next() {
		var stock Stock

		err := rows.Scan(
			&stock.ID,
			&stock.PortfolioID,
			&stock.Symbol,
			&stock.Shares,
			&stock.AveragePrice,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		p.Stocks = append(p.Stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &p, nil

}

func (ps *PortfolioStore) UpdatePortfolio(ctx context.Context, portfolioID int64, userID int64, name string) (*Portfolio, error) {
	query := `
		UPDATE portfolios 
		SET name = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, name, created_at, updated_at
	`

	var portfolio Portfolio

	err := ps.db.QueryRowContext(ctx, query, name, portfolioID, userID).Scan(
		&portfolio.ID,
		&portfolio.UserID,
		&portfolio.Name,
		&portfolio.CreatedAt,
		&portfolio.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	stockQuery := `
		SELECT id, portfolio_id, symbol, shares, average_price, created_at, updated_at
		FROM portfolio_stocks
		WHERE portfolio_id = $1
		ORDER BY symbol ASC
	`

	rows, err := ps.db.QueryContext(ctx, stockQuery, portfolioID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	portfolio.Stocks = []Stock{}
	for rows.Next() {
		var stock Stock
		err := rows.Scan(
			&stock.ID,
			&stock.PortfolioID,
			&stock.Symbol,
			&stock.Shares,
			&stock.AveragePrice,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		portfolio.Stocks = append(portfolio.Stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &portfolio, nil
}

func (ps *PortfolioStore) DeletePortfolio(ctx context.Context, portfolioID int64, userID int64) error {
	query := `
		DELETE * FROM portfolios
		WHERE id = $1 AND user_id = $2
	`

	result, err := ps.db.ExecContext(ctx, query, portfolioID, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil

}
