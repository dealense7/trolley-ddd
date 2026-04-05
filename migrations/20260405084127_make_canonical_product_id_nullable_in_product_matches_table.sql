-- +goose Up
-- +goose StatementBegin
ALTER TABLE product_matches DROP FOREIGN KEY product_matches_ibfk_2;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE product_matches MODIFY canonical_product_id BIGINT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE product_matches
    ADD CONSTRAINT product_matches_ibfk_2
        FOREIGN KEY (canonical_product_id)
            REFERENCES products(id)
            ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product_matches DROP FOREIGN KEY product_matches_ibfk_2;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE product_matches MODIFY canonical_product_id BIGINT NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE product_matches
    ADD CONSTRAINT product_matches_ibfk_2
        FOREIGN KEY (canonical_product_id)
            REFERENCES products(id);
-- +goose StatementEnd