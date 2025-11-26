CREATE TABLE IF NOT EXISTS portfolio_stocks(
    id bigserial PRIMARY KEY,
    portfolio_id bigint NOT NULL,
    symbol varchar(20) NOT NULL,
    shares numeric(20, 8) NOT NULL,
    average_price numeric(20, 8) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    FOREIGN KEY (portfolio_id) REFERENCES portfolios(id) ON DELETE CASCADE,
    UNIQUE(portfolio_id, symbol)
);