-- +goose Up
-- +goose StatementBegin
CREATE TABLE stores (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    name_local TEXT,                    -- Georgian name: კარფური
    slug TEXT NOT NULL UNIQUE,          -- URL-friendly: carrefour-ge
    logo_url TEXT,
    primary_color TEXT,                 -- Brand color for UI
    active BOOLEAN DEFAULT 1,
    scraper_config TEXT,                -- JSON config for scraper
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX idx_stores_slug ON stores(slug);
CREATE INDEX idx_stores_active ON stores(active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
