-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_matches (
    id BIGINT PRIMARY KEY,

    -- The match
    scraped_product_id BIGINT NOT NULL,
    canonical_product_id BIGINT NOT NULL,

    -- Match details
    match_type TEXT NOT NULL,  -- 'exact', 'fuzzy', 'manual', 'ml'
    confidence_score REAL NOT NULL,  -- 0.0 to 1.0

    matched_on TEXT,  -- 'barcode', 'name', 'sku', 'combined'
    match_evidence TEXT,  -- {"barcode_match": true, "name_similarity": 0.95, ...}

    -- Status
    status TEXT DEFAULT 'active',  -- 'active', 'rejected', 'superseded'

    created_by TEXT,  -- 'auto', 'admin_user_id', 'ml_model_v1'
    verified_at DATETIME,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (scraped_product_id) REFERENCES scraped_products(id) ON DELETE CASCADE,
    FOREIGN KEY (canonical_product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_matches_scraped ON product_matches(scraped_product_id);
CREATE INDEX idx_matches_canonical ON product_matches(canonical_product_id);
CREATE INDEX idx_matches_status ON product_matches(status);
CREATE INDEX idx_matches_confidence ON product_matches(confidence_score);

-- Unique constraint: one active match per scraped product
CREATE UNIQUE INDEX idx_matches_unique_active ON product_matches(scraped_product_id) WHERE status = 'active';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_matches;
-- +goose StatementEnd
