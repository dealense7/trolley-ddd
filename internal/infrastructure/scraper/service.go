package scraper

import (
	"context"
	"fmt"
	"sync"

	"github.com/dealense7/go-rates-ddd/internal/common/logger"
	"github.com/dealense7/go-rates-ddd/internal/domain/product"
	"github.com/dealense7/go-rates-ddd/internal/domain/scraper"
	"github.com/dealense7/go-rates-ddd/internal/domain/store"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/elastic"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/embedder"
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
func (s *ParserService) ScrapeAndPrint(
	context context.Context,
	repo product.Repository,
	target store.Branch,
	embederClient *embedder.Client,
	elClient *elastic.Client,
) error {
	_ = elClient.CreateIndex(context)

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

	var wg sync.WaitGroup
	maxGoroutines := 7
	// Create a buffered channel to act as a semaphore
	guard := make(chan struct{}, maxGoroutines)

	for _, p := range *products {
		wg.Add(1)
		go func(p scraper.ScrapedProduct) {
			guard <- struct{}{}
			scraped := product.NewScraped(p.ExternalID, target.ID, p.Name, p.Description, p.ImageURL)

			id, created, err := repo.InsertOrUpdate(context, scraped)
			if err != nil {
				s.log.Error("Error inserting product", zap.Error(err))
			}

			price := product.NewPrice(*id, p.Price, p.OriginalPrice, target.Country.CurrencyCode)

			if created {
				image, err := s.downloadImage(*id, p.ImageURL)
				if err != nil {
					s.log.Error("--- Image Not Downloaded ---", zap.Error(err))
				}

				//	embeddings, err := embederClient.EmbedFused(image.ImageURL(), p.Name)
				//	if err != nil {
				//		s.log.Error("--- EmbedFused ---", zap.Error(err))
				//		return err
				//	}
				//
				//	elasticProduct := elastic.Product{ProductID: fmt.Sprintf("%d", id), Embedding: embeddings}
				//
				//	err = elClient.IndexProduct(context, elasticProduct)
				//	if err != nil {
				//		s.log.Error("--- Elastic Index ---", zap.Error(err))
				//		return err
				//	}
				//
				//	matches, _ := elClient.FindSimilar(context, embeddings, 3, 0.92)
				//	if len(matches) > 0 {
				//		match := matches[0]
				//		s.log.Info("duplicate found",
				//			zap.Int64p("id", id),
				//			zap.String("match_id", match.ProductID),
				//			zap.Float64("score", match.Score))
				//	}
				//
				if err == nil {
					err = repo.AttachImageToProduct(context, *image)
					if err != nil {
						s.log.Error("--- Image not attached ---", zap.Error(err))
					}
				}
			}

			err = repo.InsertPrice(context, *price)
			if err != nil {
				s.log.Error("--- Price not added ---", zap.Error(err))
			}

			<-guard
			wg.Done()
		}(p)
	}
	wg.Wait()

	s.log.Info("--- END BATCH ---", zap.Int("total_items", len(*products)))

	return nil
}
