-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_variants (
    id BIGINT PRIMARY KEY,
    parent_product_id BIGINT NOT NULL,  -- The "master" product
    variant_product_id BIGINT NOT NULL,  -- The variant

    -- What makes this a variant?
    variant_type TEXT NOT NULL,  -- 'size', 'flavor', 'color', 'bundle', 'package'

    -- Variant details
    variant_value TEXT,  -- e.g., "2L" for size variant

    created_at DATETIME NOT NULL,

    FOREIGN KEY (parent_product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (variant_product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE(parent_product_id, variant_product_id)
);

CREATE INDEX idx_variants_parent ON product_variants(parent_product_id);
CREATE INDEX idx_variants_variant ON product_variants(variant_product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_variants;
-- +goose StatementEnd
