-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_bets (
    id                      BIGSERIAL       PRIMARY KEY,
    bet_uuid                UUID            NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    bet_no                  VARCHAR(50)     NOT NULL UNIQUE,
    session_id              BIGINT          NOT NULL ,
    ticket_id BIGINT          NOT NULL ,
    member_id               INTEGER         NOT NULL REFERENCES tbl_members(id),
    fish_type_id            INTEGER,
    cannon_type_id          SMALLINT        NOT NULL,
    elapsed_seconds         NUMERIC(18,6)   NOT NULL DEFAULT 0,
    bet_amount              NUMERIC(18,6)   NOT NULL DEFAULT 0,
    kill_reward_odd         NUMERIC(18,6)   NOT NULL DEFAULT 0,
    kill_reward             NUMERIC(18,6)   NOT NULL DEFAULT 0,
    miss_reward             NUMERIC(18,6)   NOT NULL DEFAULT 0,
    miss_reward_odd             NUMERIC(18,6)   NOT NULL DEFAULT 0,
    base_kill_rate          NUMERIC(18,6)   NOT NULL DEFAULT 0,
    rtp_adjustment_factor   NUMERIC(18,6)   NOT NULL DEFAULT 1,
    kill_rate               NUMERIC(18,6)   NOT NULL DEFAULT 0,
    random_roll             NUMERIC(18,6)   NOT NULL DEFAULT 0,
    is_kill                 BOOLEAN         NOT NULL DEFAULT FALSE,
    payout_amount           NUMERIC(18,6)   NOT NULL DEFAULT 0,
    rtp_before              NUMERIC(10,8)   NOT NULL DEFAULT 0,
    rtp_after               NUMERIC(10,8)   NOT NULL DEFAULT 0,
    target_rtp_snapshot     NUMERIC(10,8)   NOT NULL DEFAULT 0,
    -- jackpot_contribution records how much of this bet was allocated to the
    -- jackpot pool (bet_amount × effective contribution_rate at time of bet).
    -- Stored per-row so pool reconciliation does not depend on the current rate.
    jackpot_contribution    NUMERIC(18,6)   NOT NULL DEFAULT 0,
    jackpot_triggered       BOOLEAN         NOT NULL DEFAULT FALSE,
    jackpot_payout_amount   NUMERIC(18,6)   NOT NULL DEFAULT 0,
    status_id               SMALLINT        NOT NULL DEFAULT 1,
    "order"                 INTEGER,
    created_at              TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by              INTEGER, 
    updated_at              TIMESTAMPTZ,
    updated_by              INTEGER,
    deleted_at              TIMESTAMPTZ,
    deleted_by              INTEGER,
    CONSTRAINT chk_tbl_bets_bet_amount CHECK (bet_amount >= 0),
    CONSTRAINT chk_tbl_bets_reward_values CHECK (kill_reward_odd >= 0 AND kill_reward >= 0 AND miss_reward >= 0),
    CONSTRAINT chk_tbl_bets_probability_values CHECK (
        base_kill_rate >= 0
        AND base_kill_rate <= 1
        AND kill_rate >= 0
        AND kill_rate <= 1
        AND random_roll >= 0
        AND random_roll <= 1
    ),
    CONSTRAINT chk_tbl_bets_rtp_values CHECK (
        rtp_before >= 0
        AND rtp_after >= 0
        AND target_rtp_snapshot >= 0
        AND target_rtp_snapshot <= 1
    ),
    CONSTRAINT chk_tbl_bets_payout_values CHECK (
        payout_amount >= 0
        AND jackpot_contribution >= 0
        AND jackpot_payout_amount >= 0
    )
);

CREATE INDEX IF NOT EXISTS idx_tbl_bets_session_id
    ON tbl_bets (session_id, id);

-- Supports member-level bet history queries
CREATE INDEX IF NOT EXISTS idx_tbl_bets_member_id
    ON tbl_bets (member_id, id DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_bets CASCADE;
-- +goose StatementEnd
