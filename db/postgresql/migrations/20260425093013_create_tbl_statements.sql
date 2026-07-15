-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_statements (
    id BIGSERIAL PRIMARY KEY,
    statement_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    session_id BIGINT NOT NULL ,
    session_no VARCHAR(50), 
    ticket_id BIGINT REFERENCES  tbl_tickets(id),
    ticket_no VARCHAR(50),
    member_id INTEGER NOT NULL REFERENCES tbl_members(id),
    total_bet_amount NUMERIC(18,6) NOT NULL DEFAULT 0,
    total_bet_invalid_amount NUMERIC(18,6) NOT NULL DEFAULT 0,
    is_kill BOOLEAN NOT NULL DEFAULT FALSE,
    payout_amount NUMERIC(18,6) NOT NULL DEFAULT 0,
    jackpot_win_amount NUMERIC(18,6) NOT NULL DEFAULT 0,
    miss_reward             NUMERIC(18,6)   NOT NULL DEFAULT 0,   
    sync_id INTEGER,
    is_synced BOOLEAN NOT NULL DEFAULT FALSE,
    statement_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMPTZ,
    updated_by INTEGER,
    deleted_at TIMESTAMPTZ,
    deleted_by INTEGER,
    UNIQUE (statement_uuid)
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_statements;

-- +goose StatementEnd
