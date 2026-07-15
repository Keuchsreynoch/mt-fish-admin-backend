-- +goose Up
-- +goose StatementBegin
CREATE TABLE tbl_fish_render_families (
    id                  SERIAL PRIMARY KEY,
    render_family_name  VARCHAR(200)  NOT NULL,
    render_type_id      INTEGER       REFERENCES tbl_fish_render_types(id),
    spine_json_path     TEXT,
    spine_atlas_path    TEXT,
    spine_png_path      TEXT,

    -- audit columns
    status_id   SMALLINT    NOT NULL DEFAULT 1,
    "order"     INTEGER,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  INTEGER,
    updated_at  TIMESTAMPTZ,
    updated_by  INTEGER,
    deleted_at  TIMESTAMPTZ,
    deleted_by  INTEGER
);

INSERT INTO tbl_fish_render_families (id, render_family_name, render_type_id, spine_json_path, spine_atlas_path, spine_png_path) VALUES
    (1, 'normal_g', 1, NULL, NULL, NULL),
    (2, 'normal_c', 1, NULL, NULL, NULL),
    (3, 'normal_li', 1, NULL, NULL, NULL),
    (4, 'normal_t', 1, NULL, NULL, NULL),
    (5, 'normal_fsg', 1, NULL, NULL, NULL),
    (6, 'normal_l', 1, NULL, NULL, NULL),
    (7, 'normal_m', 1, NULL, NULL, NULL),
    (8, 'puffer_purple', 1, NULL, NULL, NULL),
    (9, 'puffer_yellow', 1, NULL, NULL, NULL),
    (10, 'squid_purple', 1, NULL, NULL, NULL),
    (11, 'squid_yellow', 1, NULL, NULL, NULL),
    (12, 'lobster_red', 1, NULL, NULL, NULL),
    (13, 'lobster_yellow', 1, NULL, NULL, NULL),
    (14, 'spine_crystal_crab', 2, '/fish/fish-all-star/resources/spines/crab/crab.json', '/fish/fish-all-star/resources/spines/crab/crab.atlas.txt', '/fish/fish-all-star/resources/spines/crab/crab.webp'),
    (15, 'spine_turtle', 2, '/fish/fish-all-star/resources/spines/dt/dt.json', '/fish/fish-all-star/resources/spines/dt/dt.atlas.txt', '/fish/fish-all-star/resources/spines/dt/dt.webp'),
    (16, 'spine_octopus', 2, '/fish/fish-all-star/resources/spines/boss_octopus/boss_octopus.json', '/fish/fish-all-star/resources/spines/boss_octopus/boss_octopus.atlas.txt', '/fish/fish-all-star/resources/spines/boss_octopus/boss_octopus.webp'),
    (17, 'spine_phoenix', 2, '/fish/fish-all-star/resources/spines/phoenix/phoenix.json', '/fish/fish-all-star/resources/spines/phoenix/phoenix.atlas.txt', '/fish/fish-all-star/resources/spines/phoenix/phoenix.webp'),
    (18, 'spine_crocodile', 2, '/fish/fish-all-star/resources/spines/cr/cr.json', '/fish/fish-all-star/resources/spines/cr/cr.atlas.txt', '/fish/fish-all-star/resources/spines/cr/cr.webp'),
    (19, 'spine_naga', 2, '/fish/fish-all-star/resources/spines/naga/naga.json', '/fish/fish-all-star/resources/spines/naga/naga.atlas.txt', '/fish/fish-all-star/resources/spines/naga/naga.webp');
SELECT setval('tbl_fish_render_families_id_seq', (SELECT MAX(id) FROM tbl_fish_render_families));


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_fish_render_families;
-- +goose StatementEnd
