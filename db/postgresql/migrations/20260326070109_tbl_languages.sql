-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_languages (
    id SERIAL PRIMARY KEY,
    language_name VARCHAR(50) NOT NULL,
    language_code VARCHAR(50) NOT NULL UNIQUE,
    flag VARCHAR(50),
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    updated_at TIMESTAMP,
    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_languages_code ON tbl_languages(language_code);

INSERT INTO tbl_languages (language_name, language_code, flag, status_id)
VALUES
('English', 'en', '🇺🇸', 1),
('Khmer', 'km', '🇰🇭', 1),
('Thai', 'th', '🇹🇭', 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_languages CASCADE;

-- +goose StatementEnd
