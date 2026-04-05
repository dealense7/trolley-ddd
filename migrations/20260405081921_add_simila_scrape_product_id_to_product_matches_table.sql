-- +goose Up
-- +goose StatementBegin
ALTER TABLE product_matches
    ADD COLUMN similar_scraped_product_id BIGINT,
ADD CONSTRAINT fk_similar_scraped_product
FOREIGN KEY (similar_scraped_product_id)
REFERENCES scraped_products(id)
ON
DELETE
CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product_matches
DROP
CONSTRAINT fk_similar_scraped_product,
DROP
COLUMN similar_scraped_product_id;
-- +goose StatementEnd
