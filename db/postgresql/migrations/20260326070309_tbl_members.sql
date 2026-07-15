-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_members (
    id SERIAL PRIMARY KEY,

    user_uuid UUID NOT NULL DEFAULT gen_random_uuid(),

    user_name VARCHAR(100) NOT NULL,
    login_id VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL, 

    phone_number VARCHAR(100) NOT NULL,
    profile_photo VARCHAR(255), 

    language_id INT REFERENCES tbl_languages(id),
    currency_id INT NOT NULL REFERENCES tbl_currencies(id),

    remark TEXT,
    nickname VARCHAR(100),

    login_session VARCHAR(255),

    last_login_at TIMESTAMPTZ,
    is_online BOOLEAN DEFAULT FALSE,

    status_id INT NOT NULL DEFAULT 1,
    timezone VARCHAR(50),
    pattern VARCHAR(50),

    "order" INT,

    created_by INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_by INT,
    updated_at TIMESTAMPTZ,

    deleted_by INT,
    deleted_at TIMESTAMPTZ
);

--  Seed data
INSERT INTO tbl_members (
    user_name,
    login_id,
    password,
    language_id,
    currency_id,
    remark,
    nickname,
    login_session,
    last_login_at,
    is_online,
    status_id,
    timezone,
    pattern,
    phone_number,
    profile_photo,
    "order",
    created_by,
    created_at
) VALUES
('ITTESTAA1', 'ittestaa1', '123456', 1, 1, 'Remark 1', 'IT1', 'session_1', NOW(), TRUE, 1, 'Asia/Phnom_Penh', 'pattern1', '+85510000001', 'user1.png', 1, 1, NOW()),
('ITTESTAA2', 'ittestaa2', '123456', 1, 1, 'Remark 2', 'IT2', 'session_2', NOW(), FALSE, 1, 'Asia/Phnom_Penh', 'pattern2', '+85510000002', 'user2.png', 2, 1, NOW()),
('ITTESTAA3', 'ittestaa3', '123456', 1, 1, 'Remark 3', 'IT3', 'session_3', NOW(), TRUE, 1, 'Asia/Phnom_Penh', 'pattern3', '+85510000003', 'user3.png', 3, 1, NOW()),
('ITTESTAA4', 'ittestaa4', '123456', 1, 1, 'Remark 4', 'IT2', 'session_4', NOW(), FALSE, 1, 'Asia/Phnom_Penh', 'pattern2', '+85510000002', 'user2.png', 2, 1, NOW()),
('ITTESTAA5', 'ittestaa5', '123456', 1, 1, 'Remark 5', 'IT3', 'session_5', NOW(), TRUE, 1, 'Asia/Phnom_Penh', 'pattern3', '+85510000003', 'user3.png', 3, 1, NOW());

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_members CASCADE;

-- +goose StatementEnd
