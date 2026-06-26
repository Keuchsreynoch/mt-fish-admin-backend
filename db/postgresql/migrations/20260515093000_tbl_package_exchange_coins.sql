-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_package_exchange_coins (
    id SERIAL PRIMARY KEY,

    package_name VARCHAR(100) NOT NULL,
    package_description VARCHAR(250) NOT NULL,

    currency_id INTEGER NOT NULL
        REFERENCES tbl_currencies(id),

    price_amount NUMERIC(18,6) NOT NULL DEFAULT 0,
    coin_amount NUMERIC(18,6) NOT NULL DEFAULT 0,

    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,

    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_by INTEGER,
    updated_at TIMESTAMP,

    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_package_exchange_currency_id
ON tbl_package_exchange_coins(currency_id);

INSERT INTO tbl_package_exchange_coins (
    package_name,
    package_description,
    currency_id,
    price_amount,
    coin_amount,
    status_id,
    "order",
    created_by,
    created_at
)
SELECT *
FROM (
VALUES

-- =========================
-- USD Packages
-- 1 USD = 200 coins
-- =========================

('Starter Pack USD', 'Starter package for USD top-up', 2, 5.000000, 1000.000000, 1, 1, 1, CURRENT_TIMESTAMP),
('Value Pack USD',   'Value package for USD top-up',   2, 10.000000, 2000.000000, 1, 2, 1, CURRENT_TIMESTAMP),
('Pro Pack USD',     'Pro package for USD top-up',     2, 25.000000, 5000.000000, 1, 3, 1, CURRENT_TIMESTAMP),

-- =========================
-- KHR Packages
-- 1 coin = 20 KHR
-- =========================

('Starter Pack KHR', 'Starter package for KHR top-up', 1, 20000.000000, 1000.000000, 1, 4, 1, CURRENT_TIMESTAMP),
('Value Pack KHR',   'Value package for KHR top-up',   1, 40000.000000, 2000.000000, 1, 5, 1, CURRENT_TIMESTAMP),
('Pro Pack KHR',     'Pro package for KHR top-up',     1, 100000.000000, 5000.000000, 1, 6, 1, CURRENT_TIMESTAMP),

-- =========================
-- CNY Packages
-- Example: 1 CNY = 550 KHR
-- =========================

('Starter Pack CNY', 'Starter package for CNY top-up', 3, 55.000000, 1512.500000, 1, 7, 1, CURRENT_TIMESTAMP),
('Value Pack CNY',   'Value package for CNY top-up',   3, 110.000000, 3025.000000, 1, 8, 1, CURRENT_TIMESTAMP),
('Pro Pack CNY',     'Pro package for CNY top-up',     3, 275.000000, 7562.500000, 1, 9, 1, CURRENT_TIMESTAMP),

-- =========================
-- VND Packages
-- Example: 1 VND = 0.16 KHR
-- =========================

('Starter Pack VND', 'Starter package for VND top-up', 4, 100000.000000, 800.000000, 1, 10, 1, CURRENT_TIMESTAMP),
('Value Pack VND',   'Value package for VND top-up',   4, 250000.000000, 2000.000000, 1, 11, 1, CURRENT_TIMESTAMP),
('Pro Pack VND',     'Pro package for VND top-up',     4, 500000.000000, 4000.000000, 1, 12, 1, CURRENT_TIMESTAMP)
) AS seed_data (
    package_name,
    package_description,
    currency_id,
    price_amount,
    coin_amount,
    status_id,
    "order",
    created_by,
    created_at
)
WHERE NOT EXISTS (
    SELECT 1
    FROM tbl_package_exchange_coins
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_package_exchange_coins;

-- +goose StatementEnd
