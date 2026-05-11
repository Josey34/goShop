package sqlite_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const schema = `
CREATE TABLE IF NOT EXISTS products (
    id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT NOT NULL DEFAULT '',
    price INTEGER NOT NULL DEFAULT 0, stock INTEGER NOT NULL DEFAULT 0,
    image_url TEXT NOT NULL DEFAULT '', created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS customers (
    id TEXT PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL, street TEXT NOT NULL, city TEXT NOT NULL,
    province TEXT NOT NULL, postal_code TEXT NOT NULL, password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS orders (
    id TEXT PRIMARY KEY, customer_id TEXT NOT NULL, total INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'PENDING', created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS order_items (
    id TEXT PRIMARY KEY, order_id TEXT NOT NULL, product_id TEXT NOT NULL,
    name TEXT NOT NULL, price INTEGER NOT NULL DEFAULT 0, quantity INTEGER NOT NULL DEFAULT 0
);
`

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("run schema: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
