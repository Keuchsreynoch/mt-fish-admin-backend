-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_members_coins_default_rate (
    id SERIAL PRIMARY KEY,
    currency_id INTEGER,
    rate NUMERIC(18,6),
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    updated_at TIMESTAMP,
    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

INSERT INTO tbl_members_coins_default_rate (
    currency_id,
    rate,
    status_id,
    "order",
    created_by,
    created_at
) VALUES
(1, 20.000000, 1, 1, 1, CURRENT_TIMESTAMP);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_members_coins_default_rate;

-- +goose StatementEnd