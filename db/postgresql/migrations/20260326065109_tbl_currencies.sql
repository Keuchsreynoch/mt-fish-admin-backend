-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_currencies (
    id SERIAL PRIMARY KEY,
    currency_name VARCHAR(50) NOT NULL,
    currency_code VARCHAR(50) NOT NULL UNIQUE,
    currency_symbol VARCHAR(10) NOT NULL,
    image VARCHAR(100),
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    updated_at TIMESTAMP,
    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_currencies_code ON tbl_currencies(currency_code);

INSERT INTO tbl_currencies (
    currency_name,
    currency_code,
    currency_symbol,
    status_id
)
VALUES
('Riel', 'KHR', '៛', 1),
('Dollar', 'USD', '$', 1),
('Yuan', 'CNY', '¥', 1),
('Dong', 'VND', '₫', 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_currencies CASCADE;

-- +goose StatementEnd
