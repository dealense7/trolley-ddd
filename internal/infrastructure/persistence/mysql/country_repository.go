package mysql

import (
	"context"

	"github.com/dealense7/go-rates-ddd/internal/domain/country"
	"github.com/jmoiron/sqlx"
)

type CountryRepository struct {
	db *sqlx.DB
}

func NewCountryRepo(db *sqlx.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) Insert(ctx context.Context, c *country.Country) error {
	query := `
        INSERT INTO countries (code, name, name_local, currency_code, currency_symbol, timezone)
        VALUES (:code, :name, :name_local, :currency_code, :currency_symbol, :timezone)
        ON DUPLICATE KEY UPDATE code = code;
    `
	_, err := r.db.NamedExecContext(ctx, query, c)
	return err
}
