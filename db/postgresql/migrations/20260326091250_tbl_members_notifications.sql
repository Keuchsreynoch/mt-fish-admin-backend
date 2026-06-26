-- +goose Up
-- +goose StatementBegin
CREATE TABLE  tbl_members_notifications (
    id                   SERIAL PRIMARY KEY,
    member_id            INT       NOT NULL REFERENCES tbl_members(id),
    notification_type_id INT       NOT NULL REFERENCES tbl_members_notifications_types(id),
    context              VARCHAR(100),
    subject              VARCHAR(255) NOT NULL,
    description          TEXT,
    status_id            SMALLINT  NOT NULL DEFAULT 0, -- 0: new, 1: read, 2: deleted
    "order"              INT,
    created_by           INT       NOT NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by           INT,
    updated_at           TIMESTAMP,
    deleted_by           INT,
    deleted_at           TIMESTAMP
);

-- seed
INSERT INTO tbl_members_notifications (
    member_id, notification_type_id, context, subject, description,
    status_id, "order", created_by
) VALUES
(1, 1, 'deposit',       'Deposit Successful',     'Your deposit of $100 has been processed.',   0, 1, 1),
(1, 2, 'withdraw',      'Withdrawal Processed',   'Your withdrawal of $50 has been completed.', 0, 2, 1),
(1, 3, 'exchange_coin', 'Coin Exchange Complete', 'You exchanged 500 coins successfully.',       0, 3, 1),
(1, 4, 'info',          'System Update Notice',   'The system will be updated on Sunday 2AM.',  0, 4, 1),
(1, 5, 'gift',          'You received a gift!',   'A gift of 100 bonus coins has been added.',  0, 5, 1),
(1, 6, 'maintenance',   'Scheduled Maintenance',  'Maintenance window: Sunday 2AM - 4AM UTC.',  0, 6, 1);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tbl_members_notifications;

-- +goose StatementEnd
