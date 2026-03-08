package seed

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func SeedData(ctx context.Context, db *sqlx.DB, log *zap.Logger) {
	seedCountries(ctx, db, log)
}
