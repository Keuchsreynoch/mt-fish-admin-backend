-- +goose Up
CREATE TABLE IF NOT EXISTS tbl_currency_default_rates (
    id BIGSERIAL PRIMARY KEY,
    currency_default_rate_uuid UUID,
    currency_id BIGINT NOT NULL,
    rate NUMERIC(20,6) NOT NULL,
    status_id BIGINT DEFAULT 1,
    "order" BIGINT NOT NULL DEFAULT 1,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_by BIGINT,
    updated_at TIMESTAMPTZ,
    deleted_by BIGINT,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_currency_default_rates_status_id
ON tbl_currency_default_rates(status_id);

CREATE INDEX IF NOT EXISTS idx_currency_default_rates_currency_id
ON tbl_currency_default_rates(currency_id);

CREATE INDEX IF NOT EXISTS idx_currency_default_rates_rate
ON tbl_currency_default_rates(rate);

INSERT INTO tbl_currency_default_rates (
    currency_default_rate_uuid,
    currency_id,
    rate,
    status_id,
    "order",
    created_by,
    created_at
)
VALUES
('a2a27ec1-aa9c-41b8-8644-3f330c03cdcb',1,1,1,1,1,'2024-08-02T16:05:30.301612345Z'),
('4ce2f032-8501-40dd-a91d-ee12537ca523',2,4000,1,1,1,'2024-08-02T16:05:30.301612345Z'),
('7f9d9b7e-f9b8-4b71-ae70-849c9737f5e3',3,0.150,1,1,1,'2024-08-02T16:05:30.301612345Z'),
('d7b3dab6-6e92-4671-8313-f7fbeac91215',4,100,1,1,1,'2024-08-02T16:05:30.301612345Z'),
('0bbfa220-e651-4a9e-9977-8c6d7ec6cc79',5,500,1,1,1,'2024-08-02T16:05:30.301612345Z');

-- +goose Down
DROP TABLE IF EXISTS tbl_currency_default_rates;