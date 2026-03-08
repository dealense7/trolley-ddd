package seed

import (
	"context"
	"time"

	"github.com/dealense7/go-rates-ddd/internal/domain/country"
	"go.uber.org/zap"
)

type CountrySeeder struct {
	Seeder
	Repo country.Repository
}

func (s *CountrySeeder) Run(ctx context.Context, log *zap.Logger) {
	items := []country.Country{
		*country.NewCountry("Spain", "ES", "España", "Europe/Madrid", country.CurrencyCodeEUR, country.CurrencySymbolEUR),
		*country.NewCountry("Italy", "IT", "Italia", "Europe/Rome", country.CurrencyCodeEUR, country.CurrencySymbolEUR),
		*country.NewCountry("Georgia", "GE", "საქართველო", "Asia/Tbilisi", country.CurrencyCodeGEL, country.CurrencySymbolGEL),
	}

	log.Info("Start | Seeding Countries")
	startTime := time.Now()
	for _, item := range items {

		err := s.Repo.Insert(ctx, &item)
		if err != nil {
			log.Error(err.Error())
		}
	}
	duration := time.Since(startTime).Milliseconds()
	log.Info("End | Seeding Countries /s", zap.Int64("duration_ms", duration))
}
