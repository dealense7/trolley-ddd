-- +goose Up
-- +goose StatementBegin
CREATE TABLE countries (
    id BIGINT PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,      -- ISO 3166-1 alpha-2 (GE, US, UK, etc.)
    name TEXT NOT NULL,
    name_local TEXT,                -- Local language name
    currency_code TEXT NOT NULL,    -- ISO 4217 (GEL, USD, EUR, etc.)
    currency_symbol TEXT NOT NULL,  -- ₾, $, €, etc.
    timezone TEXT NOT NULL,         -- Europe/Tbilisi
    active BOOLEAN DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_countries_code ON countries(code);
CREATE INDEX idx_countries_active ON countries(active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE countries;
-- +goose StatementEnd
