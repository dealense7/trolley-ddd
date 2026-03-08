package seed

import (
	"context"

	"go.uber.org/zap"
)

type Seeder interface {
	Run(ctx context.Context, log *zap.Logger)
}
