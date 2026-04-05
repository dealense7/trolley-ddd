package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type SearchResult struct {
	ProductID  string
	Embeddings []float64
	Score      float64
}

func (c *Client) GetSingle(ctx context.Context, id int64) (*SearchResult, error) {
	res, err := c.es.Get(IndexName, strconv.FormatInt(id, 10), c.es.Get.WithContext(ctx))

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("get: %s", b)
	}

	var raw struct {
		Source struct {
			ProductID string    `json:"product_id"`
			Embedding []float64 `json:"embedding"`
		} `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &SearchResult{
		ProductID:  raw.Source.ProductID,
		Embeddings: raw.Source.Embedding,
	}, nil
}

func (c *Client) FindSimilar(
	ctx context.Context,
	queryVector []float64,
	branchId int64,
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
			"filter": map[string]any{
				"bool": map[string]any{
					"must_not": []map[string]any{
						{
							"term": map[string]any{
								"branch_id": strconv.FormatInt(branchId, 10),
							},
						},
					},
				},
			},
		},
		"_source": []string{"product_id"},
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
