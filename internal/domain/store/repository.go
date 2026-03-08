package store

import "context"

type Repository interface {
	Insert(ctx context.Context, s *Store) (*int64, error)
	AddBranch(ctx context.Context, storeId int64, b *Branch) error
}
