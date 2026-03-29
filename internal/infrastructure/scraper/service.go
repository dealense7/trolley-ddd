package scraper

import (
	"context"
	"fmt"

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
func (s *ParserService) ScrapeAndPrint(context context.Context, repo product.Repository, target store.Branch) error {
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
		scraped := product.NewScraped(p.ExternalID, target.ID, p.Name, p.Description, p.ImageURL)

		id, err := repo.InsertOrUpdate(context, scraped)
		if err != nil {
			return err
		}

		product.NewPrice(*id, p.Price, p.OriginalPrice, target.Country.CurrencyCode)
	}

	s.log.Info("--- END BATCH ---", zap.Int("total_items", len(*products)))

	return nil
}
