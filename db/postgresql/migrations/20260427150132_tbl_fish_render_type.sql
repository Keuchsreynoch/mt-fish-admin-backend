-- +goose Up
-- +goose StatementBegin
CREATE TABLE tbl_fish_render_types (
    id          SERIAL PRIMARY KEY,
    type_code   VARCHAR(100) NOT NULL UNIQUE,

    -- audit columns
    status_id   SMALLINT     NOT NULL DEFAULT 1,
    "order"     INTEGER,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  INTEGER,
    updated_at  TIMESTAMPTZ,
    updated_by  INTEGER,
    deleted_at  TIMESTAMPTZ,
    deleted_by  INTEGER
);

INSERT INTO tbl_fish_render_types (id, type_code) VALUES
    (1, 'atlas_sprite_anim'),
    (2, 'spine');
SELECT setval('tbl_fish_render_types_id_seq', (SELECT MAX(id) FROM tbl_fish_render_types));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_fish_render_types;
-- +goose StatementEnd
