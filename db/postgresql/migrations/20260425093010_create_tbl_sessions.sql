-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_sessions (
    id BIGSERIAL PRIMARY KEY,
    session_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    session_no VARCHAR(50) NOT NULL UNIQUE,
    member_id INTEGER NOT NULL REFERENCES tbl_members(id),
    game_config_id INTEGER NOT NULL REFERENCES tbl_game_configs(id),
    started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    total_elapsed_seconds NUMERIC(18,6) NOT NULL DEFAULT 0,
    turnover NUMERIC(20,6) NOT NULL DEFAULT 0,
    payout NUMERIC(20,6) NOT NULL DEFAULT 0,
    profit NUMERIC(20,6) NOT NULL DEFAULT 0,

    rtp_current NUMERIC(10,6) NOT NULL DEFAULT 0,
    runtime_state_json JSONB NOT NULL DEFAULT '{}'::JSONB,
    current_context_index INTEGER,
    current_context_no INTEGER,
    current_group_id INTEGER,
    current_scene_id VARCHAR(100),
    current_boss_fish_type_id INTEGER,
    current_path_version_id BIGINT,
    context_started_at TIMESTAMPTZ,
    context_expires_at TIMESTAMPTZ,
    pending_next_context_index INTEGER,
    pending_next_scene_id VARCHAR(100),
    boss_scene_active BOOLEAN NOT NULL DEFAULT FALSE,
    boss_scene_lock_id VARCHAR(100),
    spawn_cursor INTEGER NOT NULL DEFAULT 0,
    snapshot_version BIGINT NOT NULL DEFAULT 0,
    client_snapshot_at TIMESTAMPTZ,
    last_bet_at TIMESTAMPTZ,
    next_allowed_bet_at TIMESTAMPTZ,
    device_meta_json JSONB NOT NULL DEFAULT '{}'::JSONB,
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMPTZ,
    updated_by INTEGER,
    deleted_at TIMESTAMPTZ,
    deleted_by INTEGER,
    UNIQUE (session_uuid)
);

CREATE INDEX IF NOT EXISTS idx_tbl_sessions_member_id
    ON tbl_sessions (member_id, status_id);

CREATE INDEX IF NOT EXISTS idx_tbl_sessions_member_active_runtime
    ON tbl_sessions (member_id, status_id, expires_at);

CREATE INDEX IF NOT EXISTS idx_tbl_sessions_runtime_group
    ON tbl_sessions (current_path_version_id, current_group_id, status_id);

CREATE INDEX IF NOT EXISTS idx_tbl_sessions_runtime_scene
    ON tbl_sessions (current_scene_id, boss_scene_active, status_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_sessions;
-- +goose StatementEnd
