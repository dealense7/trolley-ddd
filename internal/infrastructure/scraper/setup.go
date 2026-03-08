package scraper

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

func NewCollector(allowDomain []string) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(allowDomain...),
		colly.AllowURLRevisit(),
	)

	// Looks like different browsers
	extensions.RandomUserAgent(c)

	// Looks like you came from Google
	extensions.Referer(c)

	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	return c
}
