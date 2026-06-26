-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tbl_menus (
    id SERIAL PRIMARY KEY,
    menu_uuid UUID DEFAULT gen_random_uuid(),

    name VARCHAR(100) NOT NULL,
    icon VARCHAR(100) NOT NULL,
    path VARCHAR(255),

    parent_id INTEGER NOT NULL DEFAULT 0,
    status_id SMALLINT NOT NULL DEFAULT 1,
    "order" INTEGER NOT NULL DEFAULT 1,

    created_by INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_by INTEGER,
    updated_at TIMESTAMP,

    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_admin_menus_parent_id
    ON tbl_menus(parent_id);

CREATE INDEX IF NOT EXISTS idx_admin_menus_status_id
    ON tbl_menus(status_id);

CREATE INDEX IF NOT EXISTS idx_admin_menus_deleted_at
    ON tbl_menus(deleted_at);

-- Root Menus
INSERT INTO tbl_menus (
    name,
    icon,
    path,
    parent_id,
    status_id,
    "order"
)
VALUES
(
    'Dashboard',
    'mdi-view-dashboard-outline',
    '/',
    0,
    1,
    1
),
(
    'Fish Info',
    'mdi-fish',
    '/fish-info',
    0,
    1,
    2
),
(
    'Members',
    'mdi-account-group-outline',
    '/members',
    0,
    1,
    3
),
(
    'Transactions',
    'mdi-swap-horizontal',
    '/transactions',
    0,
    1,
    4
),
(
    'Statements',
    'mdi-trophy-outline',
    '/statements',
    0,
    1,
    5
),
(
    'Jackpot',
    'mdi-treasure-chest-outline',
    '/jackpot',
    0,
    1,
    6
),
(
    'Analytics',
    'mdi-chart-areaspline',
    '/analytics',
    0,
    1,
    7
),
(
    'Game Config',
    'mdi-tune',
    '/game-config',
    0,
    1,
    8
),
(
    'Agents',
    'mdi-account-tie-outline',
    '/agents',
    0,
    1,
    9
),
(
    'Settings',
    'mdi-cog-outline',
    '/settings',
    0,
    1,
    10
);

-- Transaction Children
INSERT INTO tbl_menus (
    name,
    icon,
    path,
    parent_id,
    status_id,
    "order"
)
SELECT
    'Coin',
    'mdi-bitcoin',
    '/transactions/coin',
    id,
    1,
    1
FROM tbl_menus
WHERE name = 'Transactions';

INSERT INTO tbl_menus (
    name,
    icon,
    path,
    parent_id,
    status_id,
    "order"
)
SELECT
    'Balance',
    'mdi-wallet-outline',
    '/transactions/balance',
    id,
    1,
    2
FROM tbl_menus
WHERE name = 'Transactions';

-- Jackpot Children
INSERT INTO tbl_menus (
    name,
    icon,
    path,
    parent_id,
    status_id,
    "order"
)
SELECT
    'Jackpot Pool',
    'mdi-treasure-chest',
    '/jackpot/pool',
    id,
    1,
    1
FROM tbl_menus
WHERE name = 'Jackpot';

INSERT INTO tbl_menus (
    name,
    icon,
    path,
    parent_id,
    status_id,
    "order"
)
SELECT
    'Jackpot Member Bonus',
    'mdi-account-cash-outline',
    '/jackpot/member-bonus',
    id,
    1,
    2
FROM tbl_menus
WHERE name = 'Jackpot';

INSERT INTO tbl_menus (
    name,
    icon,
    path,
    parent_id,
    status_id,
    "order"
)
SELECT
    'Jackpot History',
    'mdi-history',
    '/jackpot/history',
    id,
    1,
    3
FROM tbl_menus
WHERE name = 'Jackpot';

INSERT INTO tbl_menus (
    name,
    icon,
    path,
    parent_id,
    status_id,
    "order"
)
SELECT
    'Jackpot Ledger',
    'mdi-book-open-variant',
    '/jackpot/ledger',
    id,
    1,
    4
FROM tbl_menus
WHERE name = 'Jackpot';

INSERT INTO tbl_menus (
    name,
    icon,
    path,
    parent_id,
    status_id,
    "order"
)
SELECT
    'Jackpot Company Pop Up',
    'mdi-office-building-cog-outline',
    '/jackpot/company-pop-up',
    id,
    1,
    5
FROM tbl_menus
WHERE name = 'Jackpot';


-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_menus CASCADE;

-- +goose StatementEnd