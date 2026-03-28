package main

import (
	"context"

	"github.com/dealense7/go-rates-ddd/internal/common/cfg"
	"github.com/dealense7/go-rates-ddd/internal/common/logger"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/persistence/mysql"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/scraper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SeedersParams struct {
	fx.In

	LoggerFactory logger.Factory
}

func main() {
	app := fx.New(
		fx.Provide(
			cfg.NewConfig,
			logger.NewFactory,
			mysql.NewDB,
			mysql.ProvideRepositories,
		),

		fx.Provide(
			scraper.NewParserService,
			scraper.NewGlovoStrategy,
		),

		fx.Invoke(func(s *scraper.ParserService, glovo *scraper.GlovoStrategy) {
			s.AddStrategy(glovo)
		}),

		fx.Invoke(startParsing),
		fx.Invoke(func(shutdown fx.Shutdowner) { _ = shutdown.Shutdown() }),
	)

	app.Run()
}

func startParsing(s *scraper.ParserService, p SeedersParams, repos mysql.ReposContainer) {
	log := p.LoggerFactory.For(logger.General)
	log.Info("Seeding data started")

	ctx := context.Background()
	branches, err := repos.StoreRepo.GetAllBranches(ctx)
	if err != nil {
		log.Error("Failed to get all branches", zap.Error(err))
	}

	for _, b := range branches {
		if b.Active == false {
			continue
		}
		err = s.ScrapeAndPrint(ctx, repos.ProductRepo, b)
		if err != nil {
			log.Error("Failed to scrape data", zap.Error(err))
		}
		break
	}

	log.Info("Seeding data finished")
}
