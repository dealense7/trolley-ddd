package product

import "time"

type Translation struct {
	ID           int64
	ProductId    int64  `db:"product_id"`
	LanguageCode string `db:"language_code"`

	Name        string `db:"name"`
	Description string `db:"description"`

	NormalizedName string `db:"normalized_name"`

	CreatedAt time.Time `db:"created_at"`
}
