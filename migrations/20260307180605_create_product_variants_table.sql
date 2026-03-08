-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_variants (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    parent_product_id BIGINT NOT NULL,  -- The "master" product
    variant_product_id BIGINT NOT NULL,  -- The variant

    -- What makes this a variant?
    variant_type VARCHAR(50) NOT NULL,  -- 'size', 'flavor', 'color', 'bundle', 'package'

    -- Variant details
    variant_value VARCHAR(255),  -- e.g., "2L" for size variant

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (parent_product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (variant_product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE(parent_product_id, variant_product_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_variants_parent ON product_variants(parent_product_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_variants_variant ON product_variants(variant_product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_variants;
-- +goose StatementEnd
