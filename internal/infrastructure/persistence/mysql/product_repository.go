package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dealense7/go-rates-ddd/internal/domain/product"
	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) InsertOrUpdate(ctx context.Context, s *product.Scraped) (*int64, error) {
	var id int64

	query := `SELECT id from scraped_products where branch_id = ? and external_id = ?`
	err := r.db.GetContext(ctx, &id, query, s.BranchID, s.ExternalID)

	// If no result found create new one
	if errors.Is(err, sql.ErrNoRows) {
		query = `INSERT INTO scraped_products 
    				(branch_id, external_id, raw_name, raw_description, image_url, scraped_at) 
				 VALUES (:branch_id, :external_id, :raw_name, :raw_description, :image_url, :scraped_at)
				 `

		result, err := r.db.NamedExecContext(ctx, query, s)
		if err != nil {
			return nil, err
		}

		id, err = result.LastInsertId()
		if err != nil {
			return nil, err
		}
	} else {
		query = `UPDATE scraped_products SET scrape_count = scrape_count + 1, scraped_at = NOW() where id = :id`
		_, err = r.db.NamedExecContext(ctx, query, map[string]interface{}{"id": id})
		if err != nil {
			return nil, err
		}
	}

	return &id, nil
}
