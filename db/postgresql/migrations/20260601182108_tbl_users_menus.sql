-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_users_menus (
    id SERIAL PRIMARY KEY,
    user_menu_uuid UUID NOT NULL DEFAULT gen_random_uuid(),

    user_id INTEGER NOT NULL
        REFERENCES tbl_users(id),

    menu_id INTEGER NOT NULL
        REFERENCES tbl_menus(id),

    is_allow BOOLEAN NOT NULL DEFAULT TRUE,

    status_id SMALLINT NOT NULL DEFAULT 1,

    created_by INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_by INTEGER,
    updated_at TIMESTAMP,

    deleted_by INTEGER,
    deleted_at TIMESTAMP,

    UNIQUE(user_id, menu_id)
);

-- Admin001 full access
INSERT INTO tbl_users_menus (
    user_id,
    menu_id,
    is_allow,
    status_id,
    created_by
)
SELECT
    1,
    m.id,
    TRUE,
    1,
    1
FROM tbl_menus m
ON CONFLICT (user_id, menu_id)
DO UPDATE SET
    is_allow = TRUE,
    status_id = 1,
    deleted_by = NULL,
    deleted_at = NULL,
    updated_at = NOW();

-- Admin002 full access
INSERT INTO tbl_users_menus (
    user_id,
    menu_id,
    is_allow,
    status_id,
    created_by
)
SELECT
    2,
    m.id,
    TRUE,
    1,
    1
FROM tbl_menus m
ON CONFLICT (user_id, menu_id)
DO UPDATE SET
    is_allow = TRUE,
    status_id = 1,
    deleted_by = NULL,
    deleted_at = NULL,
    updated_at = NOW();

-- Admin003 full access
INSERT INTO tbl_users_menus (
    user_id,
    menu_id,
    is_allow,
    status_id,
    created_by
)
SELECT
    3,
    m.id,
    TRUE,
    1,
    1
FROM tbl_menus m
ON CONFLICT (user_id, menu_id)
DO UPDATE SET
    is_allow = TRUE,
    status_id = 1,
    deleted_by = NULL,
    deleted_at = NULL,
    updated_at = NOW();

-- Admin004 full access
INSERT INTO tbl_users_menus (
    user_id,
    menu_id,
    is_allow,
    status_id,
    created_by
)
SELECT
    4,
    m.id,
    TRUE,
    1,
    1
FROM tbl_menus m
ON CONFLICT (user_id, menu_id)
DO UPDATE SET
    is_allow = TRUE,
    status_id = 1,
    deleted_by = NULL,
    deleted_at = NULL,
    updated_at = NOW();

-- Admin005 full access
INSERT INTO tbl_users_menus (
    user_id,
    menu_id,
    is_allow,
    status_id,
    created_by
)
SELECT
    5,
    m.id,
    TRUE,
    1,
    1
FROM tbl_menus m
ON CONFLICT (user_id, menu_id)
DO UPDATE SET
    is_allow = TRUE,
    status_id = 1,
    deleted_by = NULL,
    deleted_at = NULL,
    updated_at = NOW();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_users_menus;

-- +goose StatementEnd