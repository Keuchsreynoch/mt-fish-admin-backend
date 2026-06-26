-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_jackpot_company_topup (
    id SERIAL PRIMARY KEY,
    jackpot_global_id INT NOT NULL DEFAULT 1,
    amount NUMERIC(20,2) NOT NULL,       
    current_amount_before NUMERIC(20,2), -- optional: jackpot pool before top-up
    current_amount_after NUMERIC(20,2),  -- optional: jackpot pool after top-up
    note TEXT,                            -- optional note
    "order" INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMP,
    updated_by INTEGER,
    deleted_at TIMESTAMP,
    deleted_by INTEGER,
    CONSTRAINT fk_tbl_jackpot_company_topup_global
        FOREIGN KEY (jackpot_global_id)
        REFERENCES tbl_jackpot_global(id),
    CONSTRAINT chk_tbl_jackpot_company_topup_amount CHECK (amount >= 0)
);

CREATE INDEX IF NOT EXISTS idx_tbl_jackpot_company_topup_global_id
    ON tbl_jackpot_company_topup (jackpot_global_id, id DESC);

CREATE INDEX IF NOT EXISTS idx_tbl_jackpot_company_topup_created_at
    ON tbl_jackpot_company_topup (created_at DESC, id DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_jackpot_company_topup;
-- +goose StatementEnd
