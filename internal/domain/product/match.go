package product

import "time"

type MatchType string

const (
	MatchTypeExact  MatchType = "exact"
	MatchTypeFuzzy  MatchType = "fuzzy"
	MatchTypeManual MatchType = "manual"
	MatchTypeML     MatchType = "ml"
)

type Match struct {
	ID int64

	ScrapedProductId   int64 `db:"scraped_product_id"`
	CanonicalProductId int64 `db:"canonical_product_id"`

	MatchType       MatchType `db:"match_type"` // 'exact', 'fuzzy', 'manual', 'ml'
	ConfidenceScore string    `db:"confidence_score"`

	MatchedOn     string `db:"matched_on"`     // 'barcode', 'name', 'sku', 'combined'
	MatchEvidence string `db:"match_evidence"` // {"barcode_match": true, "name_similarity": 0.95, ...}

	Status string `db:"status"` // Active, Rejected

	CreatedAt  time.Time `db:"created_at"`
	VerifiedAt time.Time `db:"verified_at"`
}
