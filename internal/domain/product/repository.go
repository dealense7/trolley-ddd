package product

import "context"

type Repository interface {
	InsertOrUpdate(ctx context.Context, s *Scraped) (*int64, error)
	//InsertPrice(ctx context.Context, s *Scraped) (*int64, error)
}
