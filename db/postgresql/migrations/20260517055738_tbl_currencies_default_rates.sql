-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_currencies_default_rates (
    id SERIAL PRIMARY KEY,

    default_currency_id SMALLINT NOT NULL REFERENCES tbl_currencies(id),

    currency_id SMALLINT NOT NULL REFERENCES tbl_currencies(id),

    rate NUMERIC(18,6) NOT NULL DEFAULT 0.000000,

    status_id SMALLINT NOT NULL DEFAULT 1,

    "order" INTEGER NOT NULL,

    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_by INTEGER,
    updated_at TIMESTAMP,

    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_currencies_default_rates_default_currency_id
ON tbl_currencies_default_rates(default_currency_id);

CREATE INDEX IF NOT EXISTS idx_currencies_default_rates_currency_id
ON tbl_currencies_default_rates(currency_id);

INSERT INTO tbl_currencies_default_rates (
    default_currency_id,
    currency_id,
    rate,
    status_id,
    "order"
)
SELECT *
FROM (
VALUES
(1, 1, 1.000000, 1, 1),
(1, 2, 4000.000000, 1, 2),
(1, 3, 5000, 1, 3),
(1, 4, 0.1500, 1, 4)
) AS seed_data (
    default_currency_id,
    currency_id,
    rate,
    status_id,
    "order"
)
WHERE NOT EXISTS (
    SELECT 1
    FROM tbl_currencies_default_rates
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_currencies_default_rates;

-- +goose StatementEnd
