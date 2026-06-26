-- +goose Up
-- +goose StatementBegin
CREATE TABLE tbl_fish_types (
    id                    SERIAL PRIMARY KEY,
    fish_type_name        VARCHAR(100) NOT NULL,
    runtime_kind_1based   INTEGER      NOT NULL,
    internal_kind         INTEGER      NOT NULL,
    has_runtime_config    BOOLEAN      NOT NULL DEFAULT FALSE,
    is_boss               BOOLEAN      NOT NULL DEFAULT FALSE,
    min_kill_odd          NUMERIC(10, 4),
    max_kill_odd          NUMERIC(10, 4),
    base_difficulty       NUMERIC(10, 6) NOT NULL DEFAULT 1.000000,
    difficulty_weight     NUMERIC(10, 6) NOT NULL DEFAULT 1.000000,
    allow_rtp_adjust      BOOLEAN      NOT NULL DEFAULT TRUE,
    miss_reward_enabled   BOOLEAN      NOT NULL DEFAULT FALSE,
    min_miss_reward_odd   NUMERIC(10, 4),
    max_miss_reward_odd   NUMERIC(10, 4),
    render_family_id      INTEGER      REFERENCES tbl_fish_render_families(id),
    default_state_code    VARCHAR(100),
    scale                 NUMERIC(10, 6),
    zindex                INTEGER,
    base_speed            NUMERIC(10, 6),
    walk_speed            NUMERIC(10, 6),
    behavior              VARCHAR(200),
    used_in_path          BOOLEAN      NOT NULL DEFAULT FALSE,
    path_usage_count      INTEGER      NOT NULL DEFAULT 0,
    asset_key             VARCHAR(200),
    prefab                TEXT,
    base_skin             TEXT,
    status_id             SMALLINT     NOT NULL DEFAULT 1,
    "order"               INTEGER,
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by            INTEGER,
    updated_at            TIMESTAMPTZ,
    updated_by            INTEGER,
    deleted_at            TIMESTAMPTZ,
    deleted_by            INTEGER,
    CONSTRAINT uq_tbl_fish_types_runtime_kind UNIQUE (runtime_kind_1based),
    CONSTRAINT uq_tbl_fish_types_internal_kind UNIQUE (internal_kind),
    CONSTRAINT chk_tbl_fish_types_kill_odd CHECK (
        min_kill_odd IS NULL
        OR max_kill_odd IS NULL
        OR (min_kill_odd > 0 AND max_kill_odd >= min_kill_odd)
    ),
    CONSTRAINT chk_tbl_fish_types_miss_reward_odd CHECK (
        min_miss_reward_odd IS NULL
        OR max_miss_reward_odd IS NULL
        OR (min_miss_reward_odd >= 0 AND max_miss_reward_odd >= min_miss_reward_odd)
    ),
    CONSTRAINT chk_tbl_fish_types_difficulty CHECK (
        base_difficulty > 0
        AND difficulty_weight > 0
    )
);

INSERT INTO tbl_fish_types (
    id, fish_type_name, runtime_kind_1based, internal_kind, has_runtime_config, is_boss,
    min_kill_odd, max_kill_odd,
    base_difficulty, difficulty_weight, allow_rtp_adjust,
    miss_reward_enabled, min_miss_reward_odd, max_miss_reward_odd,
    render_family_id, default_state_code, scale, zindex,
    base_speed, walk_speed, behavior,
    used_in_path, path_usage_count, asset_key, prefab, base_skin
) VALUES
--  normal fish — no miss reward below 10x
    (1,  'Green Fish',  1,  0,  TRUE, FALSE, 2,   2,   0.900000, 0.950000, TRUE, FALSE, NULL, NULL, 1,  'move', 1.1,              1,  0.25, NULL, NULL, TRUE, 3737, 'fish0',         '/fish/fish-all-star/new-collected/fish0.lh',  'resources/normal_fish/g_0001.png'),
    (2,  'Clownfish',  2,  1,  TRUE, FALSE, 3,   3,   0.950000, 0.980000, TRUE, FALSE, NULL, NULL, 2,  'move', 1.1,              2,  0.25, NULL, NULL, TRUE, 5535, 'fish1',         '/fish/fish-all-star/new-collected/fish1.lh',  'resources/normal_fish/c_0001.png'),
    (3,  'Lionfish',  3,  2,  TRUE, FALSE, 4,   4,   1.000000, 1.000000, TRUE, FALSE, NULL, NULL, 3,  'move', 1.1,              3,  0.25, NULL, NULL, TRUE, 2201, 'fish2',         '/fish/fish-all-star/new-collected/fish2.lh',  'resources/normal_fish/li_0001.png'),
    (4,  'Turtle',  4,  3,  TRUE, FALSE, 5,   5,   1.050000, 1.020000, TRUE, FALSE, NULL, NULL, 4,  'move', 1.1,              4,  0.25, NULL, NULL, TRUE, 4003, 'fish3',         '/fish/fish-all-star/new-collected/fish3.lh',  'resources/normal_fish/t_0001.png'),
    (5,  'Crawfish',  5,  4,  TRUE, FALSE, 8,   8,   1.120000, 1.050000, TRUE, FALSE, NULL, NULL, 5,  'move', 1.1,              5,  0.25, NULL, NULL, TRUE, 2953, 'fish4',         '/fish/fish-all-star/new-collected/fish4.lh',  'resources/normal_fish/fsg_0001.png'),
    (6,  'Anglerfish6',  6,  5,  TRUE, FALSE, 9,   9,   1.180000, 1.080000, TRUE, FALSE, NULL, NULL, 6,  'move', 1.1,              6,  0.25, NULL, NULL, TRUE, 2088, 'fish5',         '/fish/fish-all-star/new-collected/fish5.lh',  'resources/normal_fish/l_0001.png'),
    (7,  'Manta Ray',  7,  6,  TRUE, FALSE, 10,  10,  1.240000, 1.120000, TRUE, TRUE,  0.0500, 0.2000, 7,  'move', 1.1,              7,  0.25, NULL, NULL, TRUE, 2566, 'fish6',         '/fish/fish-all-star/new-collected/fish6.lh',  'resources/normal_fish/m_0001.png'),
--  mid fish — odds rise gradually, miss reward starts small from 10x up
    (10, 'Purple Pufferfish', 8,  7,  TRUE, FALSE, 15,  25,  1.320000, 1.150000, TRUE, TRUE,  0.1000, 0.3000, 8,  'move', 0.792,            8,  0.25, NULL, NULL, TRUE, 846,  'fish7',         '/fish/fish-all-star/new-collected/fish7.lh',  'resources/puffer/purple_anger/puffer_hong_anger_00.png'),
    (11, 'Yellow Pufferfish', 9,  8,  TRUE, FALSE, 15,  30,  1.380000, 1.180000, TRUE, TRUE,  0.1200, 0.3500, 9,  'move', 0.792,            9,  0.25, NULL, NULL, TRUE, 349,  NULL,            NULL,                                          NULL),
    (12, 'Yellow Pufferfish', 10, 9,  TRUE, FALSE, 15,  30,  1.450000, 1.200000, TRUE, TRUE,  0.1200, 0.3500, 9,  'move', 0.792,            10, 0.25, NULL, NULL, TRUE, 348,  NULL,            NULL,                                          NULL),
    (13, 'Purple Squid', 11, 10, TRUE, FALSE, 20,  30,  1.500000, 1.220000, TRUE, TRUE,  0.1500, 0.4000, 10, 'move', 0.8800000000001,  11, 0.25, NULL, NULL, TRUE, 427,  'fish10',        '/fish/fish-all-star/new-collected/fish10.lh', 'resources/squid/purple_run/0_00.png'),
    (14, 'Yellow Squid', 12, 11, TRUE, FALSE, 20,  35,  1.580000, 1.260000, TRUE, TRUE,  0.1800, 0.5000, 11, 'move', 0.8800000000001,  12, 0.25, NULL, NULL, TRUE, 249,  NULL,            NULL,                                          NULL),
    (15, 'Yellow Squid', 13, 12, TRUE, FALSE, 20,  35,  1.660000, 1.300000, TRUE, TRUE,  0.2000, 0.5500, 11, 'move', 0.8800000000001,  13, 0.25, NULL, NULL, TRUE, 279,  NULL,            NULL,                                          NULL),
    (22, 'Red Lobster', 14, 13, TRUE, FALSE, 25,  35,  1.720000, 1.320000, TRUE, TRUE,  0.2200, 0.6000, 12, 'move', 0.8800000000001,  14, 0.25, NULL, NULL, TRUE, 178,  'fish13',        '/fish/fish-all-star/new-collected/fish13.lh', 'resources/spine_fbf/lobster/red/angry_00.png'),
    (23, 'Yellow Lobster', 15, 14, TRUE, FALSE, 25,  40,  1.800000, 1.360000, TRUE, TRUE,  0.2500, 0.7000, 13, 'move', 0.8800000000001,  15, 0.25, NULL, NULL, TRUE, 70,   NULL,            NULL,                                          NULL),
    (24, 'Yellow Lobster', 16, 15, TRUE, FALSE, 25,  40,  1.900000, 1.400000, TRUE, TRUE,  0.3000, 0.8000, 13, 'move', 0.8800000000001,  16, 0.25, NULL, NULL, TRUE, 70,   NULL,            NULL,                                          NULL),
--  special fish — stronger miss reward because these feel expensive to chase
    (16, 'Crystal Crab', 17, 16, TRUE, FALSE, 30,  50,  1.950000, 1.450000, TRUE, TRUE,  0.4000, 1.0000, 14, 'run',  1.1,              17, 1.1,  50,   'path_linked_stride', TRUE, 162, 'spine_crystal_crab', NULL, NULL),
    (17, 'Turtle', 18, 17, TRUE, FALSE, 40,  60,  2.080000, 1.520000, TRUE, TRUE,  0.5000, 1.2000, 15, 'run',  0.8250000000001,  18, 1.1,  NULL, 'steady_spine_swim',  TRUE, 123, 'spine_turtle',       NULL, NULL),
    (18, 'Octopus', 19, 18, TRUE, FALSE, 60,  100, 2.250000, 1.620000, TRUE, TRUE,  0.7000, 1.6000, 16, 'run',  0.8250000000001,  19, 1.0,  NULL, 'steady_spine_swim',  TRUE, 63,  'spine_octopus',      NULL, NULL),
--  boss fish — jackpot target, very high toughness, strongest miss reward
    (19, 'Phoenix', 20, 19, TRUE, TRUE,  20,  800, 2.700000, 1.900000, FALSE, TRUE, 1.0000, 2.0000, 17, 'run',  0.8800000000001, 20, 1.0,  NULL, 'steady_spine_swim',  TRUE, 52,  'spine_phoenix',      NULL, NULL),
    (20, 'Crocodile', 21, 20, TRUE, TRUE,  50, 1000, 2.900000, 2.000000, FALSE, TRUE, 1.5000, 2.5000, 18, 'run',  0.8800000000001, 20, 1.0,  NULL, 'steady_spine_swim',  TRUE, 37,  'spine_crocodile',    NULL, NULL),
    (21, 'Naga', 22, 21, TRUE, TRUE,  50, 1200, 3.150000, 2.100000, FALSE, TRUE, 2.0000, 3.0000, 19, 'run',  1.1,             21, 1.0,  NULL, 'steady_spine_swim',  TRUE, 20,  'spine_naga',         NULL, NULL),
--  special/event fish — no odds configured yet
    (25, 'Event Fish 25', 25, 24, FALSE, FALSE, NULL, NULL, 1.000000, 1.000000, TRUE, FALSE, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
    (26, 'Event Fish 26', 26, 25, FALSE, FALSE, NULL, NULL, 1.000000, 1.000000, TRUE, FALSE, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
    (27, 'Event Fish 27', 27, 26, FALSE, FALSE, NULL, NULL, 1.000000, 1.000000, TRUE, FALSE, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
    (28, 'Event Fish 28', 28, 27, FALSE, FALSE, NULL, NULL, 1.000000, 1.000000, TRUE, FALSE, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
    (29, 'Event Fish 29', 29, 28, FALSE, FALSE, NULL, NULL, 1.000000, 1.000000, TRUE, FALSE, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
    (30, 'Event Fish 30', 30, 29, FALSE, FALSE, NULL, NULL, 1.000000, 1.000000, TRUE, FALSE, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
    (31, 'Event Fish 31', 31, 30, FALSE, FALSE, NULL, NULL, 1.000000, 1.000000, TRUE, FALSE, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL);

SELECT setval('tbl_fish_types_id_seq', (SELECT MAX(id) FROM tbl_fish_types));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_fish_types;
-- +goose StatementEnd

-- CREATE TABLE tbl_fish_types (
--     id                    SERIAL PRIMARY KEY,
--     runtime_kind_1based   INTEGER      NOT NULL,
--     internal_kind         INTEGER      NOT NULL,
--     has_runtime_config    BOOLEAN      NOT NULL DEFAULT FALSE,
--     is_boss               BOOLEAN      NOT NULL DEFAULT FALSE,
--     min_kill_odd               NUMERIC(10, 4),
--     max_kill_odd               NUMERIC(10, 4),
--     miss_reward_enabled BOOLEAN NOT NULL DEFAULT FALSE,
--     min_miss_reward_odd        NUMERIC(10, 4),
--     max_miss_reward_odd        NUMERIC(10, 4),
--     render_family_id      INTEGER      REFERENCES tbl_fish_render_families(id),
--     default_state_code    VARCHAR(100),
--     scale                 NUMERIC(10, 6),
--     zindex                INTEGER,
--     base_speed            NUMERIC(10, 6),
--     walk_speed            NUMERIC(10, 6),
--     behavior              VARCHAR(200),
--     used_in_path          BOOLEAN      NOT NULL DEFAULT FALSE,
--     path_usage_count      INTEGER      NOT NULL DEFAULT 0,
--     asset_key             VARCHAR(200),
--     prefab                TEXT,
--     base_skin             TEXT,

--     status_id   SMALLINT    NOT NULL DEFAULT 1,
--     "order"     INTEGER,
--     created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     created_by  INTEGER,
--     updated_at  TIMESTAMPTZ,
--     updated_by  INTEGER,
--     deleted_at  TIMESTAMPTZ,
--     deleted_by  INTEGER
-- );

-- INSERT INTO tbl_fish_types (id, runtime_kind_1based, internal_kind, has_runtime_config, is_boss, min_odd, max_odd, kill_rate_modifier, min_reward_odd, max_reward_odd, render_family_id, default_state_code, scale, zindex, base_speed, walk_speed, behavior, used_in_path, path_usage_count, asset_key, prefab, base_skin) VALUES
--     (1, 1, 0, TRUE, FALSE, NULL, NULL, NULL, NULL, NULL, 1, 'move', 1.1, 1, 0.25, NULL, NULL, TRUE, 3737, 'fish0', '/fish/fish-all-star/new-collected/fish0.lh', 'resources/normal_fish/g_0001.png'),
--     (2, 2, 1, TRUE, FALSE, NULL, NULL, NULL, NULL, NULL, 2, 'move', 1.1, 2, 0.25, NULL, NULL, TRUE, 5535, 'fish1', '/fish/fish-all-star/new-collected/fish1.lh', 'resources/normal_fish/c_0001.png'),
--     (3, 3, 2, TRUE, FALSE, NULL, NULL, NULL, NULL, NULL, 3, 'move', 1.1, 3, 0.25, NULL, NULL, TRUE, 2201, 'fish2', '/fish/fish-all-star/new-collected/fish2.lh', 'resources/normal_fish/li_0001.png'),
--     (4, 4, 3, TRUE, FALSE, NULL, NULL, NULL, NULL, NULL, 4, 'move', 1.1, 4, 0.25, NULL, NULL, TRUE, 4003, 'fish3', '/fish/fish-all-star/new-collected/fish3.lh', 'resources/normal_fish/t_0001.png'),
--     (5, 5, 4, TRUE, FALSE, NULL, NULL, NULL, NULL, NULL, 5, 'move', 1.1, 5, 0.25, NULL, NULL, TRUE, 2953, 'fish4', '/fish/fish-all-star/new-collected/fish4.lh', 'resources/normal_fish/fsg_0001.png'),
--     (6, 6, 5, TRUE, FALSE, NULL, NULL, NULL, NULL, NULL, 6, 'move', 1.1, 6, 0.25, NULL, NULL, TRUE, 2088, 'fish5', '/fish/fish-all-star/new-collected/fish5.lh', 'resources/normal_fish/l_0001.png'),
--     (7, 7, 6, TRUE, FALSE, NULL, NULL, NULL, NULL, NULL, 7, 'move', 1.1, 7, 0.25, NULL, NULL, TRUE, 2566, 'fish6', '/fish/fish-all-star/new-collected/fish6.lh', 'resources/normal_fish/m_0001.png'),
--     (10, 8, 7, TRUE, FALSE, 15, 60, NULL, NULL, NULL, 8, 'move', 0.792, 8, 0.25, NULL, NULL, TRUE, 846, 'fish7', '/fish/fish-all-star/new-collected/fish7.lh', 'resources/puffer/purple_anger/puffer_hong_anger_00.png'),
--     (11, 9, 8, TRUE, FALSE, 15, 125, NULL, NULL, NULL, 9, 'move', 0.792, 9, 0.25, NULL, NULL, TRUE, 349, NULL, NULL, NULL),
--     (12, 10, 9, TRUE, FALSE, 15, 210, NULL, NULL, NULL, 9, 'move', 0.792, 10, 0.25, NULL, NULL, TRUE, 348, NULL, NULL, NULL),
--     (13, 11, 10, TRUE, FALSE, 20, 75, NULL, NULL, NULL, 10, 'move', 0.8800000000000001, 11, 0.25, NULL, NULL, TRUE, 427, 'fish10', '/fish/fish-all-star/new-collected/fish10.lh', 'resources/squid/purple_run/0_00.png'),
--     (14, 12, 11, TRUE, FALSE, 20, 150, NULL, NULL, NULL, 11, 'move', 0.8800000000000001, 12, 0.25, NULL, NULL, TRUE, 249, NULL, NULL, NULL),
--     (15, 13, 12, TRUE, FALSE, 20, 245, NULL, NULL, NULL, 11, 'move', 0.8800000000000001, 13, 0.25, NULL, NULL, TRUE, 279, NULL, NULL, NULL),
--     (16, 17, 16, TRUE, FALSE, NULL, NULL, NULL, NULL, NULL, 14, 'run', 1.1, 17, 1.1, 50, 'path_linked_stride', TRUE, 162, 'spine_crystal_crab', NULL, NULL),
--     (17, 18, 17, TRUE, FALSE, NULL, NULL, NULL, NULL, NULL, 15, 'run', 0.8250000000000001, 18, 1.1, NULL, 'steady_spine_swim', TRUE, 123, 'spine_turtle', NULL, NULL),
--     (18, 19, 18, TRUE, FALSE, NULL, NULL, NULL, NULL, NULL, 16, 'run', 0.8250000000000001, 19, 1, NULL, 'steady_spine_swim', TRUE, 63, 'spine_octopus', NULL, NULL),
--     (19, 20, 19, TRUE, TRUE, NULL, NULL, NULL, NULL, NULL, 17, 'run', 0.8800000000000001, 20, 1, NULL, 'steady_spine_swim', TRUE, 52, 'spine_phoenix', NULL, NULL),
--     (20, 21, 20, TRUE, TRUE, NULL, NULL, NULL, NULL, NULL, 18, 'run', 0.8800000000000001, 20, 1, NULL, 'steady_spine_swim', TRUE, 37, 'spine_crocodile', NULL, NULL),
--     (21, 22, 21, TRUE, TRUE, NULL, NULL, NULL, NULL, NULL, 19, 'run', 1.1, 21, 1, NULL, 'steady_spine_swim', TRUE, 20, 'spine_naga', NULL, NULL),
--     (22, 14, 13, TRUE, FALSE, 25, 90, NULL, NULL, NULL, 12, 'move', 0.8800000000000001, 14, 0.25, NULL, NULL, TRUE, 178, 'fish13', '/fish/fish-all-star/new-collected/fish13.lh', 'resources/spine_fbf/lobster/red/angry_00.png'),
--     (23, 15, 14, TRUE, FALSE, 25, 175, NULL, NULL, NULL, 13, 'move', 0.8800000000000001, 15, 0.25, NULL, NULL, TRUE, 70, NULL, NULL, NULL),
--     (24, 16, 15, TRUE, FALSE, 25, 280, NULL, NULL, NULL, 13, 'move', 0.8800000000000001, 16, 0.25, NULL, NULL, TRUE, 70, NULL, NULL, NULL),
--     (25, 25, 24, FALSE, FALSE, NULL, NULL, NULL, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
--     (26, 26, 25, FALSE, FALSE, NULL, NULL, NULL, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
--     (27, 27, 26, FALSE, FALSE, NULL, NULL, NULL, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
--     (28, 28, 27, FALSE, FALSE, NULL, NULL, NULL, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
--     (29, 29, 28, FALSE, FALSE, NULL, NULL, NULL, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
--     (30, 30, 29, FALSE, FALSE, NULL, NULL, NULL, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL),
--     (31, 31, 30, FALSE, FALSE, NULL, NULL, NULL, NULL, NULL, NULL, 'move', 1, 0, 0.25, NULL, NULL, TRUE, 1, NULL, NULL, NULL);
-- SELECT setval('tbl_fish_types_id_seq', (SELECT MAX(id) FROM tbl_fish_types));



-- DROP TABLE IF EXISTS tbl_fish_types;
