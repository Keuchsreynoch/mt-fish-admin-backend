-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tbl_jackpot_member_status (
    id SERIAL PRIMARY KEY,
    member_id BIGINT NOT NULL UNIQUE,       
    is_active BOOLEAN NOT NULL DEFAULT true,  
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at TIMESTAMP,
    updated_by INTEGER,
    deleted_at TIMESTAMP,
    deleted_by INTEGER
);

INSERT INTO tbl_jackpot_member_status (
    member_id,
    is_active,
    status_id,
    "order",
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by
)
VALUES (
    1,              
    true,                -- active by default
    1,                   -- default status ID
    1,                   -- default order
    CURRENT_TIMESTAMP,   -- created_at
    1,                   -- created_by (e.g., admin ID)
    NULL,                -- updated_at
    NULL,                -- updated_by
    NULL,                -- deleted_at
    NULL                 -- deleted_by
)
ON CONFLICT (member_id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_jackpot_member_status;
-- +goose StatementEnd
