package product

import "time"

type Price struct {
	ID int64

	CanonicalProductId int64 `db:"canonical_product_id"`
	ScrapedProductId   int64 `db:"scraped_product_id"`

	Amount               string `db:"amount"`
	Currency             string `db:"currency"`
	AmountInBaseCurrency string `db:"amount_in_base_currency"`

	OriginalAmount       string `db:"original_amount"`
	DiscountedPercentage string `db:"discounted_percentage"`

	CreatedAt time.Time `db:"created_at"`
}
