package seed

import (
	"context"
	"time"

	"github.com/dealense7/go-rates-ddd/internal/domain/country"
	"github.com/dealense7/go-rates-ddd/internal/domain/store"
	"go.uber.org/zap"
)

type StoreSeeder struct {
	Seeder
	CountryRepo country.Repository
	Repo        store.Repository
}

func (s *StoreSeeder) Run(ctx context.Context, log *zap.Logger) {
	countries, _ := s.CountryRepo.GetAll(ctx)

	carrefourBranches := []store.Branch{
		*store.NewBranch(
			getCountryId("Italy", countries),
			"https://glovoapp.com/en/it/milano/stores/carrefour-milano",
			"Milano",
			store.ParseProviderGlovo,
			nil,
		),
		*store.NewBranch(
			getCountryId("Spain", countries),
			"https://glovoapp.com/en/es/madrid/stores/carrefour-madrid",
			"Madrid",
			store.ParseProviderGlovo,
			nil,
		),
		*store.NewBranch(
			getCountryId("Georgia", countries),
			"https://glovoapp.com/en/ge/tbilisi/stores/1carrefour-tbi",
			"Tbilisi",
			store.ParseProviderGlovo,
			nil,
		),
	}

	items := []store.Store{
		*store.NewStore("Carrefour", "carrefour", "carrefour.png", "#254F9B", &carrefourBranches),
	}

	log.Info("Start | Seeding Stores")
	startTime := time.Now()
	for _, item := range items {
		storeId, err := s.Repo.Insert(ctx, &item)
		if err != nil {
			log.Error(err.Error())
		}
		if storeId == nil {
			log.Error("store Id is nil")
			continue
		}

		for _, branch := range *item.Branches {
			err = s.Repo.AddBranch(ctx, *storeId, &branch)
			if err != nil {
				log.Error(err.Error())
			}
		}
	}

	duration := time.Since(startTime).Milliseconds()
	log.Info("End | Seeding Stores", zap.Int64("duration_ms", duration))
}

func getCountryId(name string, items []country.Country) int64 {
	for _, item := range items {
		if item.Name == name {
			return item.ID
		}
	}

	return 0
}
