package scraper

import (
	"time"

	"github.com/dealense7/go-rates-ddd/internal/domain/store"
)

type ScrapedProduct struct {
	ExternalID     string
	Name           string
	NormalizedName *string
	Price          int64
	OldPrice       int64
	ImageURL       string
	ScrapedAt      time.Time
}

type Strategy interface {
	Name() string
	CanParse(provider store.ParseProvider) bool
	Parse(target store.Branch) (*[]ScrapedProduct, error)
}
