package scraper

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dealense7/go-rates-ddd/internal/common/logger"
	"github.com/dealense7/go-rates-ddd/internal/domain/scraper"
	"github.com/dealense7/go-rates-ddd/internal/domain/store"
	"github.com/gocolly/colly/v2"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type GlovoStrategy struct {
	log *zap.Logger
}

var _ scraper.Strategy = (*GlovoStrategy)(nil)

func NewGlovoStrategy(logFactory logger.Factory) *GlovoStrategy {
	return &GlovoStrategy{
		log: logFactory.For(logger.Scraper).With(zap.String("strategy", "glovo")),
	}
}

func (s *GlovoStrategy) Name() string {
	return "Glovo"
}

func (s *GlovoStrategy) CanParse(provider store.ParseProvider) bool {
	return provider == store.ParseProviderGlovo
}

func (s *GlovoStrategy) Parse(target store.Branch) (*[]scraper.ScrapedProduct, error) {
	items := make([]scraper.ScrapedProduct, 0)

	// Get all category links from the page
	//links := []string{"/v4/stores/52935/addresses/326622/content/main?nodeType=DEEP_LINK&link=ortofrutta-sc.11498287/verdura-fresca-c.11497954\\"}
	links, err := s.extractCategoryLinks(target.ParseUrl)
	if err != nil {
		s.log.Error("Failed to extract category links", zap.Error(err))
		return nil, err
	}

	if len(links) == 0 {
		s.log.Warn("No category links found", zap.String("url", target.ParseUrl))
		return &items, nil
	}

	s.log.Info("Found category links", zap.Int("count", len(links)))

	maxWaitGroup := 4

	var wg sync.WaitGroup
	guard := make(chan struct{}, maxWaitGroup)

	// Fetch products from each category
	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			guard <- struct{}{}
			if err := s.fetchProducts(&items, link, target); err != nil {
				s.log.Error("Failed to fetch products",
					zap.String("link", link),
					zap.Error(err),
				)
			}

			<-guard
			wg.Done()
		}(link)
	}

	wg.Wait()

	return &items, nil
}

// extractCategoryLinks finds all category API paths from the page
func (s *GlovoStrategy) extractCategoryLinks(pageURL string) ([]string, error) {
	var scriptContent strings.Builder

	c := NewCollector([]string{"glovoapp.com", "www.glovoapp.com"})

	c.OnHTML("script", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "self.__next_f.push") {
			scriptContent.WriteString(e.Text)
		}
	})

	s.log.Info("Extracting category links", zap.String("url", pageURL))

	if err := c.Visit(pageURL); err != nil {
		return nil, fmt.Errorf("failed to visit page: %w", err)
	}

	content := scriptContent.String()
	if len(content) == 0 {
		return nil, fmt.Errorf("no script content found")
	}

	return s.parseCategoryLinks(content), nil
}

// parseCategoryLinks extracts and cleans category links from script content
func (s *GlovoStrategy) parseCategoryLinks(content string) []string {
	// Find all paths with nodeType=DEEP_LINK and -sc pattern
	re := regexp.MustCompile(`[^"]*nodeType=DEEP_LINK[^"]*-sc[^"]*`)
	matches := re.FindAllString(content, -1)

	seen := make(map[string]bool)
	links := make([]string, 0, len(matches))

	for _, match := range matches {
		cleaned := s.cleanLink(match)

		// Deduplicate
		if !seen[cleaned] {
			links = append(links, cleaned)
			seen[cleaned] = true
		}
	}

	return links
}

// cleanLink unescapes and normalizes a category link
func (s *GlovoStrategy) cleanLink(raw string) string {
	// Unescape JSON encoding (e.g., \u0026 -> &, \/ -> /)
	cleaned, err := strconv.Unquote(`"` + raw + `"`)
	if err != nil {
		// Fallback to manual replacement
		cleaned = strings.ReplaceAll(raw, `\u0026`, "&")
		cleaned = strings.ReplaceAll(cleaned, `\/`, "/")
	}

	// Ensure it starts with /
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}

	return strings.TrimSpace(cleaned)
}

// fetchProducts fetches products from a category API endpoint
func (s *GlovoStrategy) fetchProducts(items *[]scraper.ScrapedProduct, path string, target store.Branch) error {
	apiURL := s.buildAPIURL(path)

	s.log.Debug("Fetching products", zap.String("url", apiURL))

	c := NewCollector([]string{"api.glovoapp.com"})

	// Set required headers
	json := gjson.Parse(*target.ScraperConfig)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("glovo-api-version", "14")
		r.Headers.Set("glovo-app-platform", "web")
		r.Headers.Set("glovo-app-type", "customer")
		r.Headers.Set("glovo-location-city-code", json.Get("glovo-location-city-code").String())
		r.Headers.Set("glovo-location-country-code", json.Get("glovo-location-country-code").String())
	})

	c.OnResponse(func(r *colly.Response) {
		s.parseProducts(items, r.Body)
	})

	c.OnError(func(r *colly.Response, err error) {
		s.log.Error("Request failed",
			zap.String("url", r.Request.URL.String()),
			zap.Int("status", r.StatusCode),
			zap.Error(err),
		)
	})

	return c.Visit(apiURL)
}

// buildAPIURL constructs the API URL from a path
func (s *GlovoStrategy) buildAPIURL(path string) string {
	const baseURL = "https://api.glovoapp.com"

	path = strings.TrimSpace(path)

	// Normalize to v3 API
	path = strings.ReplaceAll(path, "/v4/", "/v3/")

	// Ensure v3 prefix if not present
	if !strings.HasPrefix(path, "/v3/") && !strings.HasPrefix(path, "/v4/") {
		path = "/v3" + path
	}

	// Normalize content endpoint
	path = strings.ReplaceAll(path, "/content/main", "/content")
	path = strings.ReplaceAll(path, "//", "/")

	// Clean trailing slashes and backslashes
	path = strings.TrimRight(path, "/\\")

	return baseURL + path
}

// parseProducts extracts products from API response
func (s *GlovoStrategy) parseProducts(items *[]scraper.ScrapedProduct, body []byte) {
	result := gjson.ParseBytes(body)

	// Check for API errors
	if errMsg := result.Get("error.message").String(); errMsg != "" {
		s.log.Error("API returned error", zap.String("error", errMsg))
		return
	}

	// Iterate through product items
	bodyItems := result.Get("data.body").Array()

	for _, item := range bodyItems {
		elements := item.Get("data.elements")
		if !elements.Exists() {
			continue
		}

		s.extractProducts(items, elements.Array())
	}
}

// extractProducts converts API elements to ScrapedProduct
func (s *GlovoStrategy) extractProducts(items *[]scraper.ScrapedProduct, elements []gjson.Result) {
	s.log.Info("Elements Count", zap.Int("count", len(elements)))
	for _, elem := range elements {
		data := elem.Get("data")

		// Skip items without images
		imageURL := data.Get("imageUrl").String()
		if imageURL == "" {
			continue
		}

		// Extract prices (API returns in cents/minor units)
		currentPrice := int64(data.Get("priceInfo.amount").Float() * 100)
		originalPrice := int64(data.Get("price").Float() * 100)

		product := scraper.ScrapedProduct{
			ExternalID:    data.Get("id").String(),
			Name:          data.Get("name").String(),
			Description:   data.Get("description").String(),
			Price:         currentPrice,
			OriginalPrice: originalPrice,
			ImageURL:      imageURL,
			ScrapedAt:     time.Now(),
		}

		*items = append(*items, product)
	}
}

// Helper function to validate URL
func (s *GlovoStrategy) isValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}
