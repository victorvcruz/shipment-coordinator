CREATE TABLE contracts (
                           id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           order_id       UUID NOT NULL,
                           carrier_id     UUID NOT NULL REFERENCES carriers(id) ON DELETE RESTRICT,
                           price          NUMERIC(10, 2) NOT NULL,
                           estimated_days INTEGER NOT NULL,
                           contracted_at  TIMESTAMPTZ NOT NULL,
                           created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                           updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

