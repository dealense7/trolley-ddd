-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_identifiers (
    id BIGINT PRIMARY KEY,
    product_id BIGINT NOT NULL,

    -- Identifier type
    type TEXT NOT NULL,  -- 'barcode', 'gtin', 'ean13', 'upc', 'sku', 'store_id', 'manufacturer_code'
    value TEXT NOT NULL,

    -- Confidence in this identifier
    confidence REAL DEFAULT 1.0,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE(type, value)  -- Each identifier can only map to one product
);

CREATE INDEX idx_identifiers_product ON product_identifiers(product_id);
CREATE INDEX idx_identifiers_type_value ON product_identifiers(type, value);
CREATE INDEX idx_identifiers_verified ON product_identifiers(verified);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_identifiers;
-- +goose StatementEnd
