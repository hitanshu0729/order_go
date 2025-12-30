CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    price INTEGER NOT NULL,   -- store price in smallest unit (e.g. paise/cents)
    stock INTEGER NOT NULL CHECK (stock >= 0)
);
CREATE INDEX idx_products_name ON products(name);