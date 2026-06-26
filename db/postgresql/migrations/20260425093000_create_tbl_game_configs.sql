-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_game_configs (
    id SERIAL PRIMARY KEY,
    game_code VARCHAR(30) NOT NULL UNIQUE,          -- FISH
    game_name VARCHAR(100) NOT NULL,

    rtp_target NUMERIC(6,4) NOT NULL DEFAULT 0.9500,
    rtp_floor NUMERIC(6,4) NOT NULL DEFAULT 0.8200,
    rtp_ceiling NUMERIC(6,4) NOT NULL DEFAULT 0.9800,
    jackpot_rate NUMERIC(8,6) NOT NULL DEFAULT 0.010000,  -- from turnover

    max_bullet_per_second NUMERIC(10,4) NOT NULL DEFAULT 6.0,
    allow_auto_fire BOOLEAN NOT NULL DEFAULT TRUE,
    allow_auto_lock BOOLEAN NOT NULL DEFAULT TRUE,
    kill_rate_min NUMERIC(10,6) NOT NULL DEFAULT 0.005000,
    kill_rate_max NUMERIC(10,6) NOT NULL DEFAULT 0.350000,
    rtp_adjust_strength NUMERIC(10,6) NOT NULL DEFAULT 0.350000,

    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMPTZ,
    updated_by INTEGER,
    deleted_at TIMESTAMPTZ,
    deleted_by INTEGER,
    CONSTRAINT chk_tbl_game_configs_rtp_range CHECK (
        rtp_floor >= 0
        AND rtp_target > 0
        AND rtp_ceiling <= 1
        AND rtp_floor <= rtp_target
        AND rtp_target <= rtp_ceiling
    ),
    CONSTRAINT chk_tbl_game_configs_kill_rate_range CHECK (
        kill_rate_min > 0
        AND kill_rate_max <= 1
        AND kill_rate_min <= kill_rate_max
    )
);

INSERT INTO tbl_game_configs (
    game_code,
    game_name,
    rtp_target,
    rtp_floor,
    rtp_ceiling,
    jackpot_rate,
    max_bullet_per_second,
    allow_auto_fire,
    allow_auto_lock,
    kill_rate_min,
    kill_rate_max,
    rtp_adjust_strength,
    status_id,
    "order",
    created_at,
    created_by
) VALUES (
    'FISH',
    'Fish All Star',
    0.9500,
    0.8200,
    0.9800,
    0.010000,
    6.0000,
    TRUE,
    TRUE,
    0.005000,
    0.350000,
    0.350000,
    1,
    1,
    CURRENT_TIMESTAMP,
    1
)
ON CONFLICT (game_code) DO NOTHING;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_game_configs CASCADE;

-- +goose StatementEnd
