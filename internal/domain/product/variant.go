package product

import "time"

type Variant struct {
	ID int64

	ParentProductId  int64 `db:"parent_product_id"`
	VariantProductId int64 `db:"variant_product_id"`

	VariantType  int64  `db:"variant_type"`
	VariantValue string `db:"variant_value"`

	CreatedAt time.Time `db:"created_at"`
}
