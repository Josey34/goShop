CREATE TABLE IF NOT EXISTS order_items (
    id         TEXT PRIMARY KEY,
    order_id   TEXT NOT NULL,
    product_id TEXT NOT NULL,
    name       TEXT NOT NULL,
    price      INTEGER NOT NULL DEFAULT 0,
    quantity   INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);
