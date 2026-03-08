-- +goose Up
-- +goose StatementBegin
CREATE TABLE stores (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,          -- URL-friendly: carrefour-ge
    logo_url VARCHAR(255),
    primary_color VARCHAR(255),                 -- Brand color for UI
    active BOOLEAN DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_stores_slug ON stores(slug(255));
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_stores_active ON stores(active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stores;
-- +goose StatementEnd
