package seed

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func Run(ctx context.Context, db *sqlx.DB, log *zap.Logger) {
	seedCountries(ctx, db, log)
}
