package product

import "context"

type Repository interface {
	InsertOrUpdate(ctx context.Context, s *Scraped) (*int64, bool, error)
	AttachImageToProduct(ctx context.Context, i Image) error
	InsertPrice(ctx context.Context, p Price) error
}
