package main

import (
	"context"

	"github.com/dealense7/go-rates-ddd/internal/common/cfg"
	"github.com/dealense7/go-rates-ddd/internal/common/logger"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/persistence/mysql"
	"github.com/dealense7/go-rates-ddd/scripts/seed"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// General Staff
		fx.Provide(
			cfg.NewConfig,
			logger.NewFactory,
			mysql.NewDB,
		),

		fx.Invoke(seedData),
		fx.Invoke(func(shutdown fx.Shutdowner) {
			_ = shutdown.Shutdown()
		}),
	).Run()
}

func seedData(db *sqlx.DB, loggerFactory logger.Factory) {
	log := loggerFactory.For(logger.General)
	log.Info("Seeding data started")
	ctx := context.Background()

	seed.SeedData(ctx, db, log)
}
