package country

import "context"

type Repository interface {
	GetAll(ctx context.Context) ([]Country, error)
	Insert(ctx context.Context, c *Country) error
}
