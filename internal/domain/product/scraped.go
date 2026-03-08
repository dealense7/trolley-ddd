package product

import (
	"time"
)

type MatchStatus string
type MatchMethod string

const (
	MatchStatusPending      MatchStatus = "pending"
	MatchStatusMatched      MatchStatus = "matched"
	MatchStatusNeedsReview  MatchStatus = "needs_review"
	MatchStatusNoMatch      MatchStatus = "no_match"
	MatchMethodExactBarCode MatchMethod = "exact_barcode"
	MatchMethodFuzzyName    MatchMethod = "fuzzy_name"
	MatchMethodManually     MatchMethod = "manually"
	MatchMethodML           MatchMethod = "machine_learning"
)

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

	MatchState       MatchStatus `db:"match_state"` // 'pending', 'matched', 'needs_review', 'no_match'
	MatchedProductID *int64      `db:"matched_product_id"`
	MatchConfidence  float64     `db:"match_confidence"`
	MatchMethod      MatchMethod `db:"match_method"` //'exact_barcode', 'fuzzy_name', 'manual', 'ml_model'

	ScrapedAt    time.Time `db:"scraped_at"`
	ScrapedCount int       `db:"scraped_count"`

	RawText string `db:"raw_text"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewScraped(branchId int64, rawName, image string) *Scraped {
	now := time.Now()
	return &Scraped{
		BranchID:     branchId,
		RawName:      rawName,
		ImageURL:     image,
		MatchState:   MatchStatusPending,
		ScrapedAt:    now,
		ScrapedCount: 1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (sp *Scraped) MarkAsMatched(productId int64, confidence float64, method MatchMethod) {
	sp.MatchState = MatchStatusMatched
	sp.MatchedProductID = &productId
	sp.MatchConfidence = confidence
	sp.MatchMethod = method
	sp.UpdatedAt = time.Now()
}
