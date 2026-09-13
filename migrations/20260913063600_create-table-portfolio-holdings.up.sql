CREATE TABLE portfolio_holdings (
    id BIGSERIAL PRIMARY KEY,
    portfolio_id BIGINT NOT NULL REFERENCES portfolios (id),
    stock_id BIGINT NOT NULL REFERENCES stocks (id),
    quantity BIGINT NOT NULL DEFAULT 0,
    average_buy_price NUMERIC(20, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (portfolio_id, stock_id),
    CHECK (quantity >= 0),
    CHECK (average_buy_price >= 0)
);