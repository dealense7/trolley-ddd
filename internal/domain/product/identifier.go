package product

import "time"

type Identifier struct {
	ID        int64
	ProductId int64 `db:"product_id"`

	Type  string `db:"type"`
	Value string `db:"value"`

	Confidence string `db:"confidence"`

	CreatedAt time.Time `db:"created_at"`
}
