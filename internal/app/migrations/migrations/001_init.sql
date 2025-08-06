-- +goose Up

CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS currency_prices (
     id SERIAL PRIMARY KEY,
     currency_id INTEGER REFERENCES currencies(id) ON DELETE CASCADE,
     price DECIMAL(18, 8) NOT NULL,
     timestamp BIGINT NOT NULL,
     created_at TIMESTAMP DEFAULT NOW(),

     UNIQUE(currency_id, timestamp)
);

CREATE TABLE IF NOT EXISTS scheduler_settings (
    currency_id INTEGER PRIMARY KEY REFERENCES currencies(id) ON DELETE CASCADE,
    interval_seconds INTEGER DEFAULT 60,
    last_updated_at BIGINT
);

-- +goose Down

DROP TABLE IF EXISTS scheduler_settings;
DROP TABLE IF EXISTS currency_prices;
DROP TABLE IF EXISTS currencies;

