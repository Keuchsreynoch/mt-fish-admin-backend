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


-- =====================================================
-- Root Menus
-- =====================================================

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
    'Jackpot Management',
    'mdi-treasure-chest-outline',
    '/jackpot/pool',
    0,
    1,
    2
),
(
    'Report',
    'mdi-chart-box-outline',
    '/report',
    0,
    1,
    3
),
(
    'Statement',
    'mdi-file-document-outline',
    '/statements',
    0,
    1,
    4
),
(
    'Transactions',
    'mdi-swap-horizontal',
    '/transactions',
    0,
    1,
    5
),
(
    'Game Config',
    'mdi-tune',
    '/game-config',
    0,
    1,
    6
),
(
    'Users',
    'mdi-account-group-outline',
    '/users',
    0,
    1,
    7
),
(
    'Fish Info',
    'mdi-fish',
    '/fish-info',
    0,
    1,
    8
),
(
    'History',
    'mdi-history',
    '/history',
    0,
    1,
    9
);


-- =====================================================
-- Transaction Children
-- =====================================================

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
WHERE name = 'Transactions'
  AND parent_id = 0;

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
WHERE name = 'Transactions'
  AND parent_id = 0;


-- =====================================================
-- Users Children
-- =====================================================

INSERT INTO tbl_menus (
    name,
    icon,
    path,
    parent_id,
    status_id,
    "order"
)
SELECT
    'Members',
    'mdi-account-multiple-outline',
    '/members',
    id,
    1,
    1
FROM tbl_menus
WHERE name = 'Users'
  AND parent_id = 0;

INSERT INTO tbl_menus (
    name,
    icon,
    path,
    parent_id,
    status_id,
    "order"
)
SELECT
    'Users',
    'mdi-account-outline',
    '/users',
    id,
    1,
    2
FROM tbl_menus
WHERE name = 'Users'
  AND parent_id = 0;


-- =====================================================
-- History Children
-- =====================================================

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
    1
FROM tbl_menus
WHERE name = 'History'
  AND parent_id = 0;

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
    2
FROM tbl_menus
WHERE name = 'History'
  AND parent_id = 0;

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
    3
FROM tbl_menus
WHERE name = 'History'
  AND parent_id = 0;

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
    4
FROM tbl_menus
WHERE name = 'History'
  AND parent_id = 0;


-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_menus CASCADE;

-- +goose StatementEnd