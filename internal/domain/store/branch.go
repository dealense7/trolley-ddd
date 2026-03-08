package store

import "time"

type ParseProvider string

const (
	ParseProviderGlovo ParseProvider = "glovo"
)

type Branch struct {
	ID int64 `db:"id"`

	StoreID   int64 `db:"store_id"`
	CountryId int64 `db:"country_id"`

	ParseUrl      string        `db:"parse_url"`
	ParseProvider ParseProvider `db:"parse_provider"`

	City          string  `db:"city"`
	ScraperConfig *string `db:"scraper_config"`

	Active bool `db:"active"`

	CreatedAt time.Time `db:"created_at"`
}

func NewBranch(countryId int64, parseUrl, city string, provider ParseProvider, config *string) *Branch {
	now := time.Now()
	return &Branch{
		CountryId:     countryId,
		ParseUrl:      parseUrl,
		ParseProvider: provider,
		City:          city,
		Active:        true,
		ScraperConfig: config,
		CreatedAt:     now,
	}
}
