-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_roles (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL UNIQUE,
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    updated_at TIMESTAMP,
    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

INSERT INTO tbl_roles (
    role_name,
    status_id,
    "order",
    created_by
) VALUES
    ('ADMIN', 1, 2, 1)
ON CONFLICT (role_name) DO NOTHING;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_roles;

-- +goose StatementEnd
