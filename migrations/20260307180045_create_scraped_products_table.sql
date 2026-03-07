-- +goose Up
-- +goose StatementBegin
-- Raw product data as scraped from stores
CREATE TABLE scraped_products (
    id BIGINT PRIMARY KEY,

    -- Source information
    branch_id BIGINT NOT NULL,
    external_id TEXT,  -- Store's product ID
    url TEXT,

    -- Raw product data (as-is from scraper)
    raw_name TEXT NOT NULL,
    raw_description TEXT,

    -- Identifiers found during scraping
    barcode TEXT,
    sku TEXT,
    gtin TEXT,

    -- Product details
    brand TEXT,

    -- Physical properties (as-scraped, not normalized)
    weight_value TEXT,
    weight_unit TEXT,
    volume_value TEXT,
    volume_unit TEXT,

    -- Images
    image_url TEXT,

    -- Matching status
    match_status TEXT DEFAULT 'pending',    -- 'pending', 'matched', 'needs_review', 'no_match'
    matched_product_id BIGINT,                -- Link to canonical product (if matched)
    match_confidence REAL,
    match_method TEXT,                      -- 'exact_barcode', 'fuzzy_name', 'manual', 'ml_model'

    -- Scraping metadata
    scraped_at DATETIME NOT NULL,
    scrape_count INTEGER DEFAULT 1,

    -- Raw JSON data (full scraper output)
    raw_data TEXT,

    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    FOREIGN KEY (branch_id) REFERENCES store_branches(id) ON DELETE CASCADE,
    FOREIGN KEY (matched_product_id) REFERENCES products(id) ON DELETE SET NULL
);

CREATE INDEX idx_scraped_store ON scraped_products(branch_id);
CREATE INDEX idx_scraped_barcode ON scraped_products(barcode);
CREATE INDEX idx_scraped_status ON scraped_products(match_status);
CREATE INDEX idx_scraped_matched ON scraped_products(matched_product_id);
CREATE INDEX idx_scraped_external ON scraped_products(branch_id, external_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scraped_products;
-- +goose StatementEnd
