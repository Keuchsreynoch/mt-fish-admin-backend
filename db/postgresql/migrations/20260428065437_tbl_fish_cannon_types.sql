-- +goose Up
-- +goose StatementBegin

-- Cannon levels store reusable visual assets.
-- Cannon types are exact bet denominations that point to one cannon level.

CREATE TABLE IF NOT EXISTS tbl_fish_cannon_levels (
    id               SMALLINT        PRIMARY KEY,
    level_code       VARCHAR(50)     NOT NULL UNIQUE,
    level_name       VARCHAR(100)    NOT NULL,
    cannon_frame     VARCHAR(100)    NOT NULL,
    bullet_frame     VARCHAR(100)    NOT NULL,
    net_frame        VARCHAR(100)    NOT NULL,
    status_id        SMALLINT        NOT NULL DEFAULT 1,
    "order"          SMALLINT,
    created_at       TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by       INTEGER,
    updated_at       TIMESTAMPTZ,
    updated_by       INTEGER,
    deleted_at       TIMESTAMPTZ,
    deleted_by       INTEGER
);

CREATE TABLE IF NOT EXISTS tbl_fish_cannon_types (
    id               SMALLINT        PRIMARY KEY,
    cannon_code      VARCHAR(50)     NOT NULL UNIQUE,
    cannon_name      VARCHAR(100)    NOT NULL,
    bet_amount       NUMERIC(18,6)   NOT NULL UNIQUE,
    cannon_level_id  SMALLINT        NOT NULL REFERENCES tbl_fish_cannon_levels(id),
    status_id        SMALLINT        NOT NULL DEFAULT 1,
    "order"          SMALLINT,
    created_at       TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by       INTEGER,
    updated_at       TIMESTAMPTZ,
    updated_by       INTEGER,
    deleted_at       TIMESTAMPTZ,
    deleted_by       INTEGER,

    CONSTRAINT chk_cannon_bet_amount CHECK (bet_amount > 0)
);

INSERT INTO tbl_fish_cannon_levels
    (id, level_code, level_name, cannon_frame, bullet_frame, net_frame, status_id, "order", created_by)
VALUES
    (1, 'CANNON_LEVEL_01', 'Cannon Level I',   'cannon_common01.png', 'bullet_common01.png', 'h_01.png',  1, 1, 1),
    (2, 'CANNON_LEVEL_02', 'Cannon Level II',  'cannon_common02.png', 'bullet_common02.png', 'h01_2.png', 1, 2, 1),
    (3, 'CANNON_LEVEL_03', 'Cannon Level III', 'cannon_common03.png', 'bullet_common03.png', 'h01.png',   1, 3, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO tbl_fish_cannon_types
    (id, cannon_code, cannon_name, bet_amount, cannon_level_id, status_id, "order", created_by)
VALUES
    (1,  'CANNON_BET_10',    'Cannon Bet 10',    10,    1, 1, 1, 1),
    (2,  'CANNON_BET_20',    'Cannon Bet 20',    20,    1, 1, 2, 1),
    (3,  'CANNON_BET_50',    'Cannon Bet 50',    50,    1, 1, 3, 1),
    (4,  'CANNON_BET_100',   'Cannon Bet 100',   100,   1, 1, 4, 1),
    (5,  'CANNON_BET_200',   'Cannon Bet 200',   200,   2, 1, 5, 1),
    (6,  'CANNON_BET_500',   'Cannon Bet 500',   500,   2, 1, 6, 1),
    (7,  'CANNON_BET_1000',  'Cannon Bet 1000',  1000,  2, 1, 7, 1),
    (8,  'CANNON_BET_2000',  'Cannon Bet 2000',  2000,  3, 1, 8, 1),
    (9,  'CANNON_BET_5000',  'Cannon Bet 5000',  5000,  3, 1, 9, 1),
    (10, 'CANNON_BET_10000', 'Cannon Bet 10000', 10000, 3, 1, 10, 1)
ON CONFLICT (id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_fish_cannon_types;
DROP TABLE IF EXISTS tbl_fish_cannon_levels;
-- +goose StatementEnd
