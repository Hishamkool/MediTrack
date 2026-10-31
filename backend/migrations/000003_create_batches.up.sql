-- Physical stock: one row per batch received.
CREATE TABLE batches (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    medicine_id   BIGINT      NOT NULL REFERENCES medicines(id) ON DELETE CASCADE,
    batch_number  TEXT        NOT NULL,
    expiry_date   DATE        NOT NULL,
    quantity      INTEGER     NOT NULL DEFAULT 0 CHECK (quantity >= 0), -- packs
    received_date DATE        NOT NULL DEFAULT CURRENT_DATE,
    status        TEXT        NOT NULL DEFAULT 'active'
                  CHECK (status IN ('active', 'expired', 'removed')),
    notes         TEXT        NOT NULL DEFAULT '',
    created_by    BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_batches_medicine    ON batches (medicine_id);
CREATE INDEX idx_batches_expiry_date ON batches (expiry_date);
