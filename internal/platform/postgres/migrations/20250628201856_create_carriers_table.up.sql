CREATE TABLE carriers
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE carrier_policies
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    carrier_id    UUID NOT NULL REFERENCES carriers(id) ON DELETE CASCADE,
    region        TEXT NOT NULL,
    estimated_days INTEGER NOT NULL,
    price_per_kg  NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_carrier_region UNIQUE (carrier_id, region)
);