-- +goose Up
-- +goose StatementBegin
CREATE TABLE scraped_products (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    -- Source information
    branch_id BIGINT NOT NULL,
    external_id VARCHAR(255),  -- Store's product ID
    url VARCHAR(255),

    -- Raw product data (as-is from scraper)
    raw_name VARCHAR(255) NOT NULL,
    raw_description MEDIUMTEXT,

    -- Identifiers found during scraping
    barcode VARCHAR(255),
    sku VARCHAR(255),
    gtin VARCHAR(255),

    -- Product details
    brand VARCHAR(255),

    -- Physical properties (as-scraped, not normalized)
    weight_value VARCHAR(255),
    weight_unit VARCHAR(255),
    volume_value VARCHAR(255),
    volume_unit VARCHAR(255),

    -- Images
    image_url VARCHAR(255),

    -- Matching status
    match_status VARCHAR(50) DEFAULT 'pending',    -- 'pending', 'matched', 'needs_review', 'no_match'
    matched_product_id BIGINT,                      -- Link to canonical product (if matched)
    match_confidence REAL,
    match_method VARCHAR(50),                      -- 'exact_barcode', 'fuzzy_name', 'manual', 'ml_model'

    -- Scraping metadata
    scraped_at DATETIME NOT NULL,
    scrape_count INTEGER DEFAULT 1,

    -- Raw JSON data (full scraper output)
    raw_data TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (branch_id) REFERENCES store_branches(id) ON DELETE CASCADE,
    FOREIGN KEY (matched_product_id) REFERENCES products(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_scraped_store ON scraped_products(branch_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_scraped_barcode ON scraped_products(barcode(255));
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_scraped_status ON scraped_products(match_status(50));
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_scraped_matched ON scraped_products(matched_product_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_scraped_external ON scraped_products(branch_id, external_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scraped_products;
-- +goose StatementEnd
