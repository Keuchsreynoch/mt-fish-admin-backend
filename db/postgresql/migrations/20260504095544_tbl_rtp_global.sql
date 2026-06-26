-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_rtp_global (
    id              SERIAL          PRIMARY KEY,

    -- rolling totals — updated atomically on every bet
    total_turnover  NUMERIC(30, 6)  NOT NULL DEFAULT 0,
    total_payout    NUMERIC(30, 6)  NOT NULL DEFAULT 0,

    -- pre-calculated so kill logic just reads this, no aggregation needed
    current_rtp     NUMERIC(10, 8)  NOT NULL DEFAULT 0,

    -- target your system aims to maintain
    target_rtp      NUMERIC(6, 4)   NOT NULL DEFAULT 0.8800,

    -- reset daily so RTP does not carry stale history forever
    period_date     DATE            NOT NULL DEFAULT CURRENT_DATE,

    status_id       SMALLINT        NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INTEGER,
    updated_at      TIMESTAMPTZ,
    updated_by      INTEGER,

    UNIQUE (period_date)
);

-- seed today's row
INSERT INTO tbl_rtp_global (
    total_turnover,
    total_payout,
    current_rtp,
    target_rtp,
    period_date,
    created_by
) VALUES (
    0,
    0,
    0,
    0.8800,
    CURRENT_DATE,
    1
) ON CONFLICT (period_date) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_rtp_global;
-- +goose StatementEnd