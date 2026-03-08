package store

import "context"

type Repository interface {
	GetAllBranches(ctx context.Context) ([]Branch, error)
	Insert(ctx context.Context, s *Store) (*int64, error)
	AddBranch(ctx context.Context, storeId int64, b *Branch) error
}
