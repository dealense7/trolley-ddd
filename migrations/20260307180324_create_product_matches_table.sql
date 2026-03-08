-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_matches (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    -- The match
    scraped_product_id BIGINT NOT NULL,
    canonical_product_id BIGINT NOT NULL,

    -- Match details
    match_type VARCHAR(50) NOT NULL,  -- 'exact', 'fuzzy', 'manual', 'ml'
    confidence_score REAL NOT NULL,  -- 0.0 to 1.0

    matched_on VARCHAR(50),  -- 'barcode', 'name', 'sku', 'combined'
    match_evidence VARCHAR(255),  -- {"barcode_match": true, "name_similarity": 0.95, ...}

    -- Status
    status VARCHAR(50) DEFAULT 'active',  -- 'active', 'rejected', 'superseded'

    verified_at DATETIME,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (scraped_product_id) REFERENCES scraped_products(id) ON DELETE CASCADE,
    FOREIGN KEY (canonical_product_id) REFERENCES products(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_matches_scraped ON product_matches(scraped_product_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_matches_canonical ON product_matches(canonical_product_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_matches_status ON product_matches(status(50));
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_matches_confidence ON product_matches(confidence_score);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_matches;
-- +goose StatementEnd
