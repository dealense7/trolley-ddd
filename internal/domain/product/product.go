package product

import "time"

type Product struct {
	ID int64

	CanonicalName  string `db:"canonical_name"`
	NormalizedName string `db:"normalized_name"`

	Brand           string `db:"brand"`
	CountryOfOrigin string `db:"country_of_origin"`

	NetWeight       string `db:"net_weight"`
	NetVolume       string `db:"net_volume"`
	PackageQuantity string `db:"package_quantity"`

	ImageURL string `db:"image_url"`
	Active   bool   `db:"active"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
