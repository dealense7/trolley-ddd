package main

import (
	"context"

	"github.com/dealense7/go-rates-ddd/internal/common/cfg"
	"github.com/dealense7/go-rates-ddd/internal/common/logger"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/persistence/mysql"
	"github.com/dealense7/go-rates-ddd/scripts/seed"
	"go.uber.org/fx"
)

type SeedersParams struct {
	fx.In

	LoggerFactory logger.Factory
	Seeders       []seed.Seeder `group:"seeders"`
}

func main() {
	app := fx.New(
		fx.Provide(
			cfg.NewConfig,
			logger.NewFactory,
			mysql.NewDB,
			mysql.ProvideRepositories,
			fx.Annotate(
				func(repos mysql.ReposContainer) seed.Seeder { return &seed.CountrySeeder{Repo: repos.CountryRepo} },
				fx.ResultTags(`group:"seeders"`),
			),
		),

		fx.Invoke(runAllSeeders),
		fx.Invoke(func(shutdown fx.Shutdowner) { _ = shutdown.Shutdown() }),
	)

	app.Run()
}

func runAllSeeders(p SeedersParams) {
	log := p.LoggerFactory.For(logger.General)
	log.Info("Seeding data started")

	for _, s := range p.Seeders {
		s.Run(context.Background(), log)
	}

	log.Info("Seeding data finished")
}
