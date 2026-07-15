-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_company_ledger (
    id SERIAL PRIMARY KEY,
    game_config_id INT NOT NULL REFERENCES tbl_game_configs(id),
    fish_type_id BIGINT NOT NULL REFERENCES tbl_fish_types(id),
    member_id INTEGER,
    bet_id BIGINT,
    ticket_id BIGINT,
    statement_id BIGINT,
    source_type VARCHAR(30) NOT NULL DEFAULT 'bet',
    company_profit_coin NUMERIC(20,2) NOT NULL,
    pool_before NUMERIC(20,2) NOT NULL DEFAULT 0,
    pool_after NUMERIC(20,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMPTZ,
    updated_by INTEGER,
    deleted_at TIMESTAMPTZ,
    deleted_by INTEGER,
    CONSTRAINT chk_tbl_company_ledger_amounts CHECK (
        company_profit_coin >= 0
        AND pool_before >= 0
        AND pool_after >= 0
    )
);

CREATE INDEX IF NOT EXISTS idx_tbl_company_ledger_game_config_id
    ON tbl_company_ledger (game_config_id, id DESC);

CREATE INDEX IF NOT EXISTS idx_tbl_company_ledger_fish_type_id
    ON tbl_company_ledger (fish_type_id, id DESC);

CREATE INDEX IF NOT EXISTS idx_tbl_company_ledger_created_at
    ON tbl_company_ledger (created_at DESC, id DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_company_ledger;

-- +goose StatementEnd
