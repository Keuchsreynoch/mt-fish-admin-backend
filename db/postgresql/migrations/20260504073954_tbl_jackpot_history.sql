-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_jackpot_history (
    id SERIAL PRIMARY KEY,
    jackpot_global_id INT NOT NULL DEFAULT 1,
    fish_type_id BIGINT NOT NULL,
    member_id INTEGER,
    jackpot_type SMALLINT NOT NULL DEFAULT 1,
    payout_coin NUMERIC(20,2) NOT NULL,  
    pool_before NUMERIC(20,2) NOT NULL, 
    pool_after NUMERIC(20,2) NOT NULL,   
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMP,
    updated_by INTEGER,
    deleted_at TIMESTAMP,
    deleted_by INTEGER,
    CONSTRAINT fk_tbl_jackpot_history_global
        FOREIGN KEY (jackpot_global_id)
        REFERENCES tbl_jackpot_global(id),
    CONSTRAINT fk_tbl_jackpot_history_fish_type
        FOREIGN KEY (fish_type_id)
        REFERENCES tbl_fish_types(id),
    CONSTRAINT chk_tbl_jackpot_history_amounts CHECK (
        payout_coin >= 0
        AND pool_before >= 0
        AND pool_after >= 0
    )
);

CREATE INDEX IF NOT EXISTS idx_tbl_jackpot_history_global_id
    ON tbl_jackpot_history (jackpot_global_id, id DESC);

CREATE INDEX IF NOT EXISTS idx_tbl_jackpot_history_fish_type_id
    ON tbl_jackpot_history (fish_type_id, id DESC);

CREATE INDEX IF NOT EXISTS idx_tbl_jackpot_history_created_at
    ON tbl_jackpot_history (created_at DESC, id DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_jackpot_history;
-- +goose StatementEnd
