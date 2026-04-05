CREATE TABLE IF NOT EXISTS customers (
    id            TEXT PRIMARY KEY,
    name          TEXT NOT NULL,
    email         TEXT NOT NULL UNIQUE,
    phone         TEXT NOT NULL,
    street        TEXT NOT NULL,
    city          TEXT NOT NULL,
    province      TEXT NOT NULL,
    postal_code   TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at    DATETIME NOT NULL,
    updated_at    DATETIME NOT NULL
);
