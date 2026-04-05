-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_images
(
    id             BIGINT PRIMARY KEY AUTO_INCREMENT,

    product_id     BIGINT,
    size           BIGINT,

    name           varchar(255),
    folder         varchar(255),
    extension      varchar(255),

    has_embeddings BOOLEAN            DEFAULT 0,

    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (product_id) REFERENCES scraped_products (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_product_prices_scraped ON product_images (product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_images;
-- +goose StatementEnd