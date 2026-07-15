-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_jackpot_member_bonus (
    id SERIAL PRIMARY KEY,
    member_id INTEGER NOT NULL,
    amount NUMERIC(20,2) NOT NULL,
    note TEXT,
    "order" INTEGER NOT NULL DEFAULT 1,
    status_id SMALLINT NOT NULL DEFAULT 1, -- 1 pending
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMP,
    updated_by INTEGER,
    deleted_at TIMESTAMP,
    deleted_by INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_jackpot_member_bonus;
-- +goose StatementEnd
