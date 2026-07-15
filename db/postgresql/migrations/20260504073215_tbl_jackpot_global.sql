-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_jackpot_global (
    id SERIAL PRIMARY KEY,
    current_amount NUMERIC(20,2) NOT NULL DEFAULT 0,       
    threshold_amount NUMERIC(20,2) NOT NULL DEFAULT 1000000, 
    chance_denom BIGINT NOT NULL DEFAULT 5000,
    payout_percent NUMERIC(5,2) NOT NULL DEFAULT 0, 
    min_eligible_bet_amount NUMERIC(20,2) NOT NULL DEFAULT 0,
    company_topup_amount NUMERIC(20,2) NOT NULL DEFAULT 0,     -- amount added by company to boost jackpot
    jackpot_fixed_payout_amount NUMERIC(20,2) NOT NULL DEFAULT 0, -- actual fixed payout when jackpot hits
    status_id               SMALLINT        NOT NULL DEFAULT 1,
    "order"                 INTEGER,
    created_at              TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by              INTEGER, 
    updated_at              TIMESTAMPTZ,
    updated_by              INTEGER,
    deleted_at              TIMESTAMPTZ,
    deleted_by              INTEGER,
    CONSTRAINT chk_tbl_jackpot_global_amounts CHECK (
        current_amount >= 0
        AND threshold_amount > 0
        AND payout_percent >= 0
        AND payout_percent <= 100
        AND min_eligible_bet_amount >= 0
        AND company_topup_amount >= 0
        AND jackpot_fixed_payout_amount >= 0
    ),
    CONSTRAINT chk_tbl_jackpot_global_chance CHECK (chance_denom > 0)
);

INSERT INTO tbl_jackpot_global (
    current_amount, 
    threshold_amount, 
    chance_denom, 
    payout_percent,
    min_eligible_bet_amount,
    company_topup_amount,
    jackpot_fixed_payout_amount,
    status_id,
    "order",
    created_at,
    created_by
)
VALUES (
    0,          -- starting jackpot amount
    1000000,    -- threshold to trigger jackpot
    5000,       -- chance denominator
    50,         -- payout percent
    0,          -- minimum eligible bet
    0,          -- company top-up to jackpot
    0,          -- extra payout from company
    1,
    1,
    CURRENT_TIMESTAMP,
    1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_jackpot_global;
-- +goose StatementEnd
