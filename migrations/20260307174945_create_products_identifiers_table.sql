-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_identifiers (
    id BIGINT PRIMARY KEY,
    product_id BIGINT NOT NULL,

    -- Identifier type
    type VARCHAR(50) NOT NULL,  -- 'barcode', 'gtin', 'ean13', 'upc', 'sku', 'store_id', 'manufacturer_code'
    value VARCHAR(255) NOT NULL,

    -- Confidence in this identifier
    confidence REAL DEFAULT 1.0,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE(type, value)  -- Each identifier can only map to one product
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_identifiers_product ON product_identifiers(product_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_identifiers_type_value ON product_identifiers(type(50), value(255));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_identifiers;
-- +goose StatementEnd
