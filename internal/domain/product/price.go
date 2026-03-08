package product

import (
	"time"

	"github.com/dealense7/go-rates-ddd/internal/domain/country"
)

type Price struct {
	ID int64

	CanonicalProductId *int64 `db:"canonical_product_id"`
	ScrapedProductId   int64  `db:"scraped_product_id"`

	Amount               int64                `db:"amount"`
	Currency             country.CurrencyCode `db:"currency"`
	AmountInBaseCurrency string               `db:"amount_in_base_currency"`

	OriginalAmount       string `db:"original_amount"`
	DiscountedPercentage string `db:"discounted_percentage"`

	CreatedAt time.Time `db:"created_at"`
}

func NewPrice(productId, amount int64, currency country.CurrencyCode) *Price {
	now := time.Now()
	return &Price{
		ScrapedProductId: productId,
		Amount:           amount,
		Currency:         currency,
		CreatedAt:        now,
	}
}
