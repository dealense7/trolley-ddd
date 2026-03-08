package scraper

import (
	"fmt"
	"time"

	"github.com/dealense7/go-rates-ddd/internal/common/logger"
	"github.com/dealense7/go-rates-ddd/internal/domain/product"
	"github.com/dealense7/go-rates-ddd/internal/domain/scraper"
	"github.com/dealense7/go-rates-ddd/internal/domain/store"
	"go.uber.org/zap"
)

type ParserService struct {
	log        *zap.Logger
	strategies []scraper.Strategy
}

func NewParserService(logFactory logger.Factory) *ParserService {
	return &ParserService{
		log:        logFactory.For(logger.Scraper),
		strategies: []scraper.Strategy{},
	}
}

func (s *ParserService) AddStrategy(st scraper.Strategy) {
	s.strategies = append(s.strategies, st)
}

// ScrapeAndPrint just logs the data, no saving
func (s *ParserService) ScrapeAndPrint(target store.Branch) error {
	var str scraper.Strategy
	for _, st := range s.strategies {
		if st.CanParse(target.ParseProvider) {
			str = st
			break
		}
	}
	if str == nil {
		return fmt.Errorf("no strategy found for %s", target.ParseUrl)
	}

	products, err := str.Parse(target)
	if err != nil {
		return err
	}

	for _, p := range *products {
		ExternalID     string
		Price          int64
		OldPrice       int64
		ScrapedAt      time.Time
		scraped := product.NewScraped(target.ID, p.Name, p.ImageURL)
		price := product.newP
	}

	s.log.Info("--- END BATCH ---", zap.Int("total_items", len(*products)))

	return nil
}
