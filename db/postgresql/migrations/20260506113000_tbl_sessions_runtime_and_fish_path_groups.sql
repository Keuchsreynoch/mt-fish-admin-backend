-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_fish_path_versions (
    id BIGSERIAL PRIMARY KEY,
    version_code VARCHAR(100) NOT NULL,
    version_name VARCHAR(255),
    content_hash VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    source_name VARCHAR(255),
    remark TEXT,
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMPTZ,
    updated_by INTEGER,
    deleted_at TIMESTAMPTZ,
    deleted_by INTEGER,
    CONSTRAINT uq_tbl_fish_path_versions_code UNIQUE (version_code)
);

CREATE TABLE IF NOT EXISTS tbl_fish_path_groups (
    id BIGSERIAL PRIMARY KEY,
    path_version_id BIGINT NOT NULL REFERENCES tbl_fish_path_versions(id),
    group_id INTEGER NOT NULL,
    context_no INTEGER,
    duration_seconds NUMERIC(12,3),
    has_boss BOOLEAN NOT NULL DEFAULT FALSE,
    fish_count INTEGER NOT NULL DEFAULT 0,
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMPTZ,
    updated_by INTEGER,
    deleted_at TIMESTAMPTZ,
    deleted_by INTEGER,
    CONSTRAINT uq_tbl_fish_path_groups_version_group UNIQUE (path_version_id, group_id)
);

CREATE TABLE IF NOT EXISTS tbl_fish_path_group_fish (
    id BIGSERIAL PRIMARY KEY,
    path_group_id BIGINT NOT NULL REFERENCES tbl_fish_path_groups(id) ON DELETE CASCADE,
    fish_type_id INTEGER NOT NULL REFERENCES tbl_fish_types(id),
    path_id INTEGER,
    spawn_delay_seconds NUMERIC(12,3) NOT NULL DEFAULT 0,
    is_boss BOOLEAN NOT NULL DEFAULT FALSE,
    occurrence_count INTEGER NOT NULL DEFAULT 1,
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMPTZ,
    updated_by INTEGER,
    deleted_at TIMESTAMPTZ,
    deleted_by INTEGER
);

CREATE INDEX IF NOT EXISTS idx_tbl_fish_path_versions_active
    ON tbl_fish_path_versions (is_active, status_id);

CREATE INDEX IF NOT EXISTS idx_tbl_fish_path_groups_version_context
    ON tbl_fish_path_groups (path_version_id, context_no, group_id);

CREATE INDEX IF NOT EXISTS idx_tbl_fish_path_group_fish_group
    ON tbl_fish_path_group_fish (path_group_id, fish_type_id, is_boss);

INSERT INTO tbl_fish_path_versions (
    id, version_code, version_name, content_hash, is_active, source_name, remark,
    status_id, "order", created_by, updated_by
) VALUES (
    1,
    'fish-all-star-v1',
    'Fish All Star Path V1',
    NULL,
    TRUE,
    'public/fish/fish-all-star/path.json',
    'Initial active path version for session/runtime validation.',
    1,
    1,
    1,
    1
)
ON CONFLICT (version_code) DO UPDATE
SET
    version_name = EXCLUDED.version_name,
    is_active = EXCLUDED.is_active,
    source_name = EXCLUDED.source_name,
    remark = EXCLUDED.remark,
    updated_at = CURRENT_TIMESTAMP,
    updated_by = EXCLUDED.updated_by;

SELECT setval(
    pg_get_serial_sequence('tbl_fish_path_versions', 'id'),
    GREATEST((SELECT COALESCE(MAX(id), 1) FROM tbl_fish_path_versions), 1)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_tbl_fish_path_group_fish_group;
DROP INDEX IF EXISTS idx_tbl_fish_path_groups_version_context;
DROP INDEX IF EXISTS idx_tbl_fish_path_versions_active;

DROP TABLE IF EXISTS tbl_fish_path_group_fish;
DROP TABLE IF EXISTS tbl_fish_path_groups;
DROP TABLE IF EXISTS tbl_fish_path_versions;

-- +goose StatementEnd
