-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_tickets (
    id BIGSERIAL PRIMARY KEY,
    ticket_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    ticket_no VARCHAR(50) NOT NULL UNIQUE,
    session_id BIGINT NOT NULL ,
    member_id INTEGER NOT NULL REFERENCES tbl_members(id),
    bet_amount NUMERIC(18,6) NOT NULL DEFAULT 0,
    valid_bet_amount NUMERIC(18,6) NOT NULL DEFAULT 0,
    payout_amount NUMERIC(18,6) NOT NULL DEFAULT 0,
    jackpot_win_amount NUMERIC(18,6) NOT NULL DEFAULT 0,
    miss_reward             NUMERIC(18,6)   NOT NULL DEFAULT 0,    
    platform_id INTEGER ,
    domain_name VARCHAR(50),
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMPTZ,
    updated_by INTEGER,
    deleted_at TIMESTAMPTZ,
    deleted_by INTEGER,
    UNIQUE (ticket_uuid)
);

CREATE INDEX IF NOT EXISTS idx_tbl_tickets_session_id
    ON tbl_tickets (session_id, status_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_tickets CASCADE;

-- +goose StatementEnd
