package country

import "context"

type Repository interface {
	Insert(ctx context.Context, c *Country) error
}
