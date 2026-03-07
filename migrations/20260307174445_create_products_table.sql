-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    id BIGINT PRIMARY KEY,

    canonical_name TEXT NOT NULL,   -- Primary name (usually English)
    normalized_name TEXT NOT NULL,  -- Normalized name for matching (lowercase, no special chars)

    -- Product details
    brand TEXT,
    country_of_origin TEXT,

    -- Physical properties
    net_weight REAL,           -- Converted to grams
    net_volume REAL,           -- Converted to ml
    package_quantity INTEGER,  -- Number of items in package

    -- Primary image
    image_url TEXT,

    active BOOLEAN DEFAULT 1,

    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX idx_products_normalized ON products(normalized_name);
CREATE INDEX idx_products_brand ON products(brand);
CREATE INDEX idx_products_active ON products(active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
