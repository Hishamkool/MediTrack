-- Dealers / medical reps we order from.
CREATE TABLE dealers (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name         TEXT        NOT NULL,
    phone        TEXT        NOT NULL DEFAULT '',
    alt_phone    TEXT        NOT NULL DEFAULT '',
    company_name TEXT        NOT NULL DEFAULT '',
    notes        TEXT        NOT NULL DEFAULT '',
    is_active    BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- The medicine master list.
CREATE TABLE medicines (
    id                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    composition       TEXT          NOT NULL,
    strength          TEXT          NOT NULL DEFAULT '',
    brand_name        TEXT          NOT NULL,
    manufacturer      TEXT          NOT NULL DEFAULT '',
    dosage_form       TEXT          NOT NULL DEFAULT '',
    units_per_pack    INTEGER       CHECK (units_per_pack > 0),
    default_dealer_id BIGINT        REFERENCES dealers(id) ON DELETE SET NULL,
    ptr               NUMERIC(12,2) CHECK (ptr >= 0),        -- price to retailer, per pack
    offer_remarks     TEXT          NOT NULL DEFAULT '',     -- e.g. "10+5"
    reorder_level     INTEGER       CHECK (reorder_level >= 0),
    is_active         BOOLEAN       NOT NULL DEFAULT TRUE,
    created_by        BIGINT        REFERENCES users(id) ON DELETE SET NULL,
    created_at        TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- lower(...) indexes make case-insensitive search fast.
CREATE INDEX idx_medicines_brand_name     ON medicines (lower(brand_name));
CREATE INDEX idx_medicines_composition    ON medicines (lower(composition));
CREATE INDEX idx_medicines_default_dealer ON medicines (default_dealer_id);
