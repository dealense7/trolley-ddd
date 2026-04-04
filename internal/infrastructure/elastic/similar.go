package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type SearchResult struct {
	ProductID string
	Score     float64
}

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
			Score:     h.Score,
		})
	}
	return results, nil
}
