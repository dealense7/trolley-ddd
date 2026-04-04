package elastic

import (
	"fmt"

	"github.com/dealense7/go-rates-ddd/internal/common/cfg"
	es8 "github.com/elastic/go-elasticsearch/v9"
)

const IndexName = "products"

// Vector dimensions from CLIP
const VectorDims = 512

type Client struct {
	es *es8.Client
}

func New(cfg *cfg.Config) (*Client, error) {
	ecfg := es8.Config{Addresses: []string{cfg.ElasticsearchURL}}
	c, err := es8.NewClient(ecfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elastic: %w", err)
	}

	return &Client{es: c}, nil
}

type Product struct {
	ProductID string    `json:"product_id"`
	BranchId  string    `json:"branch_id"`
	Embedding []float64 `json:"embedding"`
}
