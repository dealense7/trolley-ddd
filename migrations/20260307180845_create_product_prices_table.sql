-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_prices (
    id BIGINT PRIMARY KEY,

    -- Link to product (one or the other)
    canonical_product_id BIGINT,  -- Direct link to canonical product
    scraped_product_id BIGINT,    -- Link to scraped product (which links to canonical)

    -- Price information
    amount REAL NOT NULL,
    currency VARCHAR(10) NOT NULL,
    amount_in_base_currency REAL,

    original_amount REAL,
    discount_percentage REAL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (canonical_product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (scraped_product_id) REFERENCES scraped_products(id) ON DELETE CASCADE,

    -- Must have one or the other
    CHECK (
        (canonical_product_id IS NOT NULL AND scraped_product_id IS NULL) OR
        (canonical_product_id IS NULL AND scraped_product_id IS NOT NULL)
        )
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_product_prices_canonical ON product_prices(canonical_product_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_product_prices_scraped ON product_prices(scraped_product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_prices;
-- +goose StatementEnd
