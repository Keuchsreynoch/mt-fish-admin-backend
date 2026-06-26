-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_users (
    id BIGSERIAL PRIMARY KEY,

    user_uuid UUID NOT NULL DEFAULT gen_random_uuid(),

    user_name VARCHAR(100) NOT NULL,
    login_id VARCHAR(100) NOT NULL UNIQUE,

    email VARCHAR(50),
    password VARCHAR(255) NOT NULL,

    nickname VARCHAR(100),

    role_id INTEGER NOT NULL,
    role_name VARCHAR(50) NOT NULL,
 
    parent_id INTEGER ,

    is_admin BOOLEAN DEFAULT FALSE,

    login_session VARCHAR(255),
    last_login_at TIMESTAMPTZ,

    profile VARCHAR(100),

    status_id SMALLINT NOT NULL DEFAULT 1,

    "order" INTEGER,

    created_by INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_by INTEGER,
    updated_at TIMESTAMP,

    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

-- 🌱 SEED DATA
INSERT INTO tbl_users (
    user_name,
    login_id,
    email,
    password,
    nickname,
    role_id,
    role_name,
    is_admin,
    status_id,
    "order",
    created_by,
    created_at
)
SELECT
    v.user_name,
    v.login_id,
    v.email,
    '123456',
    v.nickname,
    1,
    'ADMIN',
    TRUE,
    1,
    v.display_order,
    1,
    NOW()
FROM (
    VALUES
        ('ADMIN001', 'ADMIN001', 'admin001@example.com', '123456', 1),
        ('ADMIN002', 'ADMIN002', 'admin002@example.com', '123456', 2),
        ('ADMIN003', 'ADMIN003', 'admin003@example.com', '123456', 3),
        ('ADMIN004', 'ADMIN004', 'admin004@example.com', '123456', 4),
        ('ADMIN005', 'ADMIN005', 'admin005@example.com', '123456', 5)
) AS v(user_name, login_id, email, nickname, display_order)
ON CONFLICT (login_id) DO NOTHING;  

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_users CASCADE;

-- +goose StatementEnd
