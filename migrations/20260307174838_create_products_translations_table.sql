-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_translations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_id BIGINT NOT NULL,
    language_code VARCHAR(3) NOT NULL,  -- ISO 639-1: en, ka, ru, de, fr, etc.

    -- Translated fields
    name  VARCHAR(255) NOT NULL,
    description  MEDIUMTEXT,

    -- Normalized version for matching
    normalized_name VARCHAR(255) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE(product_id, language_code)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_translations_product ON product_translations(product_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_translations_language ON product_translations(language_code(3));
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_translations_normalized ON product_translations(normalized_name(255));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_translations;
-- +goose StatementEnd
