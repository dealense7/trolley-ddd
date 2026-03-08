package seed

import (
	"context"
	"time"

	"github.com/dealense7/go-rates-ddd/internal/domain/country"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func seedCountries(ctx context.Context, db *sqlx.DB, log *zap.Logger) {
	items := []country.Country{
		*country.NewCountry("Spain", "ES", "España", "Europe/Madrid", country.CurrencyCodeEUR, country.CurrencySymbolEUR),
		*country.NewCountry("Italy", "IT", "Italia", "Europe/Rome", country.CurrencyCodeEUR, country.CurrencySymbolEUR),
		*country.NewCountry("Georgia", "GE", "საქართველო", "Asia/Tbilisi", country.CurrencyCodeGEL, country.CurrencySymbolGEL),
	}

	query := `
		INSERT INTO countries (code, name, name_local, currency_code, currency_symbol, timezone)
		VALUES (:code, :name, :name_local, :currency_code, :currency_symbol, :timezone)
		ON DUPLICATE KEY UPDATE code = code;
	`

	log.Info("Start | Seeding Countries")
	startTime := time.Now()
	for _, item := range items {
		_, err := db.NamedExecContext(ctx, query, item)
		if err != nil {
			log.Error(err.Error())
		}
	}
	duration := time.Since(startTime).Milliseconds()
	log.Info("End | Seeding Countries /s", zap.Int64("duration_ms", duration))

}
