-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_members_coins (
    id SERIAL PRIMARY KEY,
    member_id INTEGER NOT NULL,
    coin_date TIMESTAMP,
    beginning_coin NUMERIC(18,6) DEFAULT 0,
    deposit NUMERIC(18,6) DEFAULT 0,
    withdraw NUMERIC(18,6) DEFAULT 0,
    coin_amount NUMERIC(18,6) DEFAULT 0,
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    updated_at TIMESTAMP,
    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

INSERT INTO tbl_members_coins (
    member_id,
    coin_date,
    beginning_coin,
    deposit,
    withdraw,
    coin_amount,
    status_id,
    "order",
    created_by,
    created_at
) VALUES
(1, CURRENT_TIMESTAMP, 0.000000,   100.000000, 0.000000,   100.000000, 1, 1, 1, CURRENT_TIMESTAMP),
(2, CURRENT_TIMESTAMP, 200.000000, 50.000000,  20.000000,  230.000000, 1, 2, 1, CURRENT_TIMESTAMP),
(3, CURRENT_TIMESTAMP, 0.000000,   75.000000,  0.000000,   75.000000,  1, 3, 1, CURRENT_TIMESTAMP),
(4, CURRENT_TIMESTAMP, 500.000000, 0.000000,   100.000000, 400.000000, 1, 4, 1, CURRENT_TIMESTAMP),
(5, CURRENT_TIMESTAMP, 20.000000,  30.000000,  10.000000,  40.000000,  1, 5, 1, CURRENT_TIMESTAMP);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_members_coins;

-- +goose StatementEnd
