package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	es8 "github.com/elastic/go-elasticsearch/v9"
)

const IndexName = "products"

// Vector dimensions from BGE-Visualized-M3
const VectorDims = 1024

type Client struct {
	es *es8.Client
}

func New(addr string) (*Client, error) {
	cfg := es8.Config{Addresses: []string{addr}}
	c, err := es8.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{es: c}, nil
}

// ── Index mapping ─────────────────────────────────────────────

// CreateIndex creates the products index with dense_vector mapping.
// Call once at startup (safe to call multiple times — skips if exists).
func (c *Client) CreateIndex(ctx context.Context) error {
	mapping := map[string]any{
		"mappings": map[string]any{
			"properties": map[string]any{
				// The fused embedding vector
				"embedding": map[string]any{
					"type":       "dense_vector",
					"dims":       VectorDims,
					"index":      true,
					"similarity": "cosine", // cosine similarity for normalized vectors
				},
				// Stored metadata (not used for search, just returned)
				"product_id": map[string]any{"type": "keyword"},
				"name":       map[string]any{"type": "text"},
				"country":    map[string]any{"type": "keyword"},
				"image_url":  map[string]any{"type": "keyword", "index": false},
				"created_at": map[string]any{"type": "date"},
			},
		},
		"settings": map[string]any{
			"number_of_shards":   1,
			"number_of_replicas": 0, // single node — no replicas needed
		},
	}

	body, _ := json.Marshal(mapping)
	res, err := c.es.Indices.Create(
		IndexName,
		c.es.Indices.Create.WithContext(ctx),
		c.es.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 400 = already exists, that's fine
	if res.IsError() && res.StatusCode != 400 {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("create index: %s", b)
	}
	return nil
}

// ── Document schema ───────────────────────────────────────────

type Product struct {
	ProductID string    `json:"product_id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	ImageURL  string    `json:"image_url"`
	Embedding []float64 `json:"embedding"`
}

// ── Index (save) a product ────────────────────────────────────

func (c *Client) IndexProduct(ctx context.Context, p Product) error {
	body, err := json.Marshal(p)
	if err != nil {
		return err
	}

	res, err := c.es.Index(
		IndexName,
		bytes.NewReader(body),
		c.es.Index.WithContext(ctx),
		c.es.Index.WithDocumentID(p.ProductID), // upsert by product ID
		c.es.Index.WithRefresh("false"),        // async refresh (faster bulk indexing)
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("index product %s: %s", p.ProductID, b)
	}
	return nil
}

// ── Search ─────────────────────────────────────────────────────

type SearchResult struct {
	ProductID string
	Name      string
	Country   string
	ImageURL  string
	Score     float64
}

// FindSimilar performs kNN vector search.
// threshold: minimum cosine score (0.0–1.0). Use ~0.85 for "same product".
func (c *Client) FindSimilar(
	ctx context.Context,
	queryVector []float64,
	topK int,
	threshold float64,
) ([]SearchResult, error) {

	// Elasticsearch 8.x kNN query
	query := map[string]any{
		"knn": map[string]any{
			"field":          "embedding",
			"query_vector":   queryVector,
			"k":              topK,
			"num_candidates": topK * 10, // oversample for better recall
			"similarity":     threshold, // minimum score filter
		},
		"_source": []string{"product_id", "name", "country", "image_url"},
	}

	body, _ := json.Marshal(query)
	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(IndexName),
		c.es.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("search: %s", b)
	}

	// Parse response
	var raw struct {
		Hits struct {
			Hits []struct {
				Score  float64 `json:"_score"`
				Source struct {
					ProductID string `json:"product_id"`
					Name      string `json:"name"`
					Country   string `json:"country"`
					ImageURL  string `json:"image_url"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(raw.Hits.Hits))
	for _, h := range raw.Hits.Hits {
		results = append(results, SearchResult{
			ProductID: h.Source.ProductID,
			Name:      h.Source.Name,
			Country:   h.Source.Country,
			ImageURL:  h.Source.ImageURL,
			Score:     h.Score,
		})
	}
	return results, nil
}

// BulkIndex indexes many products efficiently in one HTTP call.
func (c *Client) BulkIndex(ctx context.Context, products []Product) error {
	var buf strings.Builder
	for _, p := range products {
		meta := fmt.Sprintf(`{"index":{"_index":%q,"_id":%q}}`, IndexName, p.ProductID)
		doc, _ := json.Marshal(p)
		buf.WriteString(meta + "\n")
		buf.WriteString(string(doc) + "\n")
	}

	res, err := c.es.Bulk(
		strings.NewReader(buf.String()),
		c.es.Bulk.WithContext(ctx),
		c.es.Bulk.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("bulk index: %s", b)
	}
	return nil
}
