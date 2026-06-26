-- +goose Up
-- +goose StatementBegin

CREATE TABLE tbl_members_notifications_types (
    id                     SERIAL PRIMARY KEY,
    notification_type_name VARCHAR(100) NOT NULL,
    icon                   VARCHAR(100),
    status_id              SMALLINT     NOT NULL DEFAULT 1,
    "order"                INT,
    created_by             INT          NOT NULL,
    created_at             TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by             INT,
    updated_at             TIMESTAMP,
    deleted_by             INT,
    deleted_at             TIMESTAMP
);

-- seed
INSERT INTO tbl_members_notifications_types 
    (notification_type_name, icon, status_id, "order", created_by)
VALUES
    ('Deposit',       'icon_deposit.png',       1, 1, 1),
    ('Withdraw',      'icon_withdraw.png',       1, 2, 1),
    ('Exchange Coin', 'icon_exchange_coin.png',  1, 3, 1),
    ('Info',          'icon_info.png',           1, 4, 1),
    ('Gift',          'icon_gift.png',           1, 5, 1),
    ('Maintenance',   'icon_maintenance.png',    1, 6, 1);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_members_notifications_types CASCADE;

-- +goose StatementEnd
