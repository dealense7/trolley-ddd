-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_translations (
    id BIGINT PRIMARY KEY,
    product_id BIGINT NOT NULL,
    language_code TEXT NOT NULL,  -- ISO 639-1: en, ka, ru, de, fr, etc.

    -- Translated fields
    name TEXT NOT NULL,
    description TEXT,

    -- Normalized version for matching
    normalized_name TEXT NOT NULL,

    created_at DATETIME NOT NULL,

    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE(product_id, language_code)
);

CREATE INDEX idx_translations_product ON product_translations(product_id);
CREATE INDEX idx_translations_language ON product_translations(language_code);
CREATE INDEX idx_translations_normalized ON product_translations(normalized_name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_translations;
-- +goose StatementEnd
