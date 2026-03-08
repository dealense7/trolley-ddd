-- +goose Up
-- +goose StatementBegin
CREATE TABLE store_branches (
    id BIGINT PRIMARY KEY,
    store_id BIGINT NOT NULL,
    country_id BIGINT NOT NULL,
    parse_url TEXT,
    parse_provider TEXT,
    city TEXT NOT NULL,
    latitude REAL,
    longitude REAL,
    opening_hours TEXT,
    active BOOLEAN DEFAULT 1,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE,
    FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE RESTRICT
);

CREATE INDEX idx_store_branches_country ON store_branches(country_id);
CREATE INDEX idx_store_branches_store ON store_branches(store_id);
CREATE INDEX idx_store_branches_city ON store_branches(city);
CREATE INDEX idx_store_branches_coords ON store_branches(latitude, longitude);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS store_branches;
-- +goose StatementEnd
