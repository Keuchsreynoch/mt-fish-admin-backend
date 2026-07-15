-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_members_coins_transactions (
    id SERIAL PRIMARY KEY,
    member_id INTEGER,
    member_coin_id INTEGER,
    before_coin NUMERIC(18,6),
    amount NUMERIC(18,6),
    currency_id NUMERIC(18,6),
    transaction_type_id INTEGER,
    transaction_group_type_id INTEGER,
    transaction_date TIMESTAMP,
    require_approval BOOLEAN DEFAULT FALSE,
    reference VARCHAR(50),
    remark VARCHAR(250),
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    updated_at TIMESTAMP,
    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

INSERT INTO tbl_members_coins_transactions (
    member_id,
    member_coin_id,
    before_coin,
    amount,
    currency_id,
    transaction_type_id,
    transaction_group_type_id,
    transaction_date,
    require_approval,
    reference,
    remark,
    status_id,
    "order",
    created_by,
    created_at
) VALUES
(1, 1, 100.000000, 50.000000, 1.000000, 1, 1, CURRENT_TIMESTAMP, FALSE, 'REF00001', 'Deposit coins',   1, 1, 1, CURRENT_TIMESTAMP),
(2, 2, 230.000000, 20.000000, 1.000000, 2, 1, CURRENT_TIMESTAMP, FALSE, 'REF00002', 'Withdraw coins',  1, 2, 1, CURRENT_TIMESTAMP),
(3, 3, 75.000000,  30.000000, 1.000000, 1, 1, CURRENT_TIMESTAMP, FALSE, 'REF00003', 'Deposit coins',   1, 3, 1, CURRENT_TIMESTAMP),
(4, 4, 400.000000, 100.000000,1.000000, 3, 2, CURRENT_TIMESTAMP, TRUE,  'REF00004', 'Bet transaction', 1, 4, 1, CURRENT_TIMESTAMP),
(5, 5, 40.000000,  25.000000, 1.000000, 4, 2, CURRENT_TIMESTAMP, FALSE, 'REF00005', 'Win transaction', 1, 5, 1, CURRENT_TIMESTAMP);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_members_coins_transactions;

-- +goose StatementEnd
