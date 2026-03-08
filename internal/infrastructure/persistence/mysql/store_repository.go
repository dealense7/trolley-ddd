package mysql

import (
	"context"

	"github.com/dealense7/go-rates-ddd/internal/domain/store"
	"github.com/jmoiron/sqlx"
)

type StoreRepository struct {
	db *sqlx.DB
}

func NewStoreRepo(db *sqlx.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) Insert(ctx context.Context, c *store.Store) (*int64, error) {
	query := `
        INSERT INTO stores (name, slug, logo_url, primary_color)
        VALUES (:name, :slug, :logo_url, :primary_color)
        ON DUPLICATE KEY UPDATE slug = slug;
    `

	result, err := r.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if id == 0 {
		err := r.db.GetContext(ctx, &id, "SELECT id FROM stores WHERE slug = ?", c.Slug)
		if err != nil {
			return nil, err
		}
	}

	return &id, nil
}

func (r *StoreRepository) AddBranch(ctx context.Context, storeId int64, b *store.Branch) error {
	query := `
		INSERT IGNORE INTO store_branches (store_id, country_id, parse_url, parse_provider, city, scraper_config)
		VALUES (:store_id, :country_id, :parse_url, :parse_provider, :city, :scraper_config);
	`

	b.StoreID = storeId

	_, err := r.db.NamedExecContext(ctx, query, b)
	if err != nil {
		return err
	}

	return nil
}
