CREATE INDEX IF NOT EXISTS idx_portfolios_user_id ON portfolios(id);
CREATE INDEX IF NOT EXISTS idx_portfolio_stocks_portfolio_id ON portfolio_stocks(portfolio_id);