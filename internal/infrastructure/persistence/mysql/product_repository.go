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

func (r *ProductRepository) InsertOrUpdate(ctx context.Context, s *product.Scraped) (*int64, bool, error) {
	var id int64
	created := false

	query := `SELECT id from scraped_products where branch_id = ? and external_id = ?`
	err := r.db.GetContext(ctx, &id, query, s.BranchID, s.ExternalID)

	// If no result found create new one
	if errors.Is(err, sql.ErrNoRows) {
		created = true
		query = `INSERT INTO scraped_products 
    				(branch_id, external_id, raw_name, raw_description, image_url, scraped_at) 
				 VALUES (:branch_id, :external_id, :raw_name, :raw_description, :image_url, :scraped_at)
				 `

		result, err := r.db.NamedExecContext(ctx, query, s)
		if err != nil {
			return nil, false, err
		}

		id, err = result.LastInsertId()
		if err != nil {
			return nil, false, err
		}
	} else {
		query = `UPDATE scraped_products SET scrape_count = scrape_count + 1, scraped_at = NOW() where id = :id`
		_, err = r.db.NamedExecContext(ctx, query, map[string]interface{}{"id": id})
		if err != nil {
			return nil, false, err
		}
	}

	return &id, created, nil
}

func (r *ProductRepository) AttachImageToProduct(ctx context.Context, i product.Image) error {

	query := `INSERT INTO product_images 
    				(product_id, name, size, extension, folder) 
				 VALUES (:product_id, :name, :size, :extension, :folder)
				 `

	_, err := r.db.NamedExecContext(ctx, query, i)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) InsertPrice(ctx context.Context, p product.Price) error {

	query := `INSERT INTO product_prices
    				(scraped_product_id, amount, currency, original_amount, created_at) 
				 VALUES (:scraped_product_id, :amount, :currency, :original_amount, :created_at)
				 `

	_, err := r.db.NamedExecContext(ctx, query, p)
	if err != nil {
		return err
	}

	return nil
}
