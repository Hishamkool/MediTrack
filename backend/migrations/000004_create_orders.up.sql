-- Purchase order requests, one dealer per order.
CREATE TABLE orders (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dealer_id    BIGINT      NOT NULL REFERENCES dealers(id) ON DELETE RESTRICT,
    status       TEXT        NOT NULL DEFAULT 'draft'
                 CHECK (status IN ('draft', 'sent', 'committed', 'received', 'cancelled')),
    committed_at TIMESTAMPTZ,            -- set when status becomes 'committed'
    notes        TEXT        NOT NULL DEFAULT '',
    created_by   BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_dealer ON orders (dealer_id);
CREATE INDEX idx_orders_status ON orders (status);

-- Lines of an order. ptr/offer are copied at order time so later
-- price changes in the medicine master don't alter old orders.
CREATE TABLE order_items (
    id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id       BIGINT        NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    medicine_id    BIGINT        NOT NULL REFERENCES medicines(id) ON DELETE RESTRICT,
    quantity       INTEGER       NOT NULL CHECK (quantity > 0),
    ptr_at_order   NUMERIC(12,2),
    offer_at_order TEXT          NOT NULL DEFAULT ''
);

CREATE INDEX idx_order_items_order ON order_items (order_id);
