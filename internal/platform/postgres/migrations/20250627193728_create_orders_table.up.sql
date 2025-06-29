CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE orders
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product        TEXT                     NOT NULL,
    weight_kg      NUMERIC(10, 2)         NOT NULL,
    destination_uf CHAR(2)                  NOT NULL,
    status         TEXT                     NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
