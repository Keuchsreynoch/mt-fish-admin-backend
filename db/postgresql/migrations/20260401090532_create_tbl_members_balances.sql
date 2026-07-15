-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_members_balances (
    id SERIAL PRIMARY KEY,

    member_id INTEGER NOT NULL
        REFERENCES tbl_members(id),

    currency_id SMALLINT NOT NULL
        REFERENCES tbl_currencies(id),

    balance NUMERIC(40,3) NOT NULL DEFAULT 0.000,

    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,

    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_by INTEGER,
    updated_at TIMESTAMP,

    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_members_balances_member_id
ON tbl_members_balances(member_id);

CREATE INDEX idx_members_balances_currency_id
ON tbl_members_balances(currency_id);

CREATE UNIQUE INDEX uq_members_balances_member_currency
ON tbl_members_balances(member_id, currency_id);

INSERT INTO tbl_members_balances (
    member_id,
    currency_id,
    balance,
    status_id,
    "order",
    created_by,
    created_at
)
VALUES

-- Member 1
(1, 1, 1000.000, 1, 1, 1, CURRENT_TIMESTAMP), -- USD
(1, 2, 4000000.000, 1, 2, 1, CURRENT_TIMESTAMP), -- KHR
(1, 3, 7000.000, 1, 3, 1, CURRENT_TIMESTAMP), -- CNY
(1, 4, 25000000.000, 1, 4, 1, CURRENT_TIMESTAMP), -- VND

-- Member 2
(2, 1, 2500.500, 1, 5, 1, CURRENT_TIMESTAMP),
(2, 2, 10000000.000, 1, 6, 1, CURRENT_TIMESTAMP),
(2, 3, 15000.000, 1, 7, 1, CURRENT_TIMESTAMP),
(2, 4, 50000000.000, 1, 8, 1, CURRENT_TIMESTAMP),

-- Member 3
(3, 1, 200000.000, 1, 9, 1, CURRENT_TIMESTAMP),
(3, 2, 30000.000, 1, 10, 1, CURRENT_TIMESTAMP),
(3, 3, 5000.000, 1, 11, 1, CURRENT_TIMESTAMP),
(3, 4, 18000000.000, 1, 12, 1, CURRENT_TIMESTAMP),

-- Member 4
(4, 1, 3200.000, 1, 13, 1, CURRENT_TIMESTAMP),
(4, 2, 12800000.000, 1, 14, 1, CURRENT_TIMESTAMP),
(4, 3, 22000.000, 1, 15, 1, CURRENT_TIMESTAMP),
(4, 4, 65000000.000, 1, 16, 1, CURRENT_TIMESTAMP),

-- Member 5
(5, 1, 500.000, 1, 17, 1, CURRENT_TIMESTAMP),
(5, 2, 2000000.000, 1, 18, 1, CURRENT_TIMESTAMP),
(5, 3, 3500.000, 1, 19, 1, CURRENT_TIMESTAMP),
(5, 4, 10000000.000, 1, 20, 1, CURRENT_TIMESTAMP);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_members_balances;

-- +goose StatementEnd