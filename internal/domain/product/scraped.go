package product

import "time"

type Scraped struct {
	ID int64

	BranchID   int64  `db:"branch_id"`
	ExternalID string `db:"external_id"`
	Url        string `db:"url"`

	RawName        string `db:"raw_name"`
	RawDescription string `db:"raw_description"`

	Barcode string `db:"barcode"`
	Sku     string `db:"sku"`
	Gtin    string `db:"gtin"`

	Brand string `db:"brand"`

	WeightValue string `db:"weight_value"`
	WeightUnit  string `db:"weight_unit"`
	VolumeValue string `db:"volume_value"`
	VolumeUnit  string `db:"volume_unit"`

	ImageURL string `db:"image_url"`

	MatchState      string `db:"match_state"` // 'pending', 'matched', 'needs_review', 'no_match'
	MatchProductID  int64  `db:"match_product_id"`
	MatchConfidence string `db:"match_confidence"`
	MatchMethod     string `db:"match_method"` //'exact_barcode', 'fuzzy_name', 'manual', 'ml_model'

	ScrapedAt    time.Time `db:"scraped_at"`
	ScrapedCount int       `db:"scraped_count"`

	RawText string `db:"raw_text"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
