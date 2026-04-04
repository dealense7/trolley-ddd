package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// Creates Table
func (c *Client) CreateIndex(ctx context.Context) error {
	mapping := map[string]any{
		"mappings": map[string]any{
			"properties": map[string]any{
				"embedding": map[string]any{
					"type":       "dense_vector",
					"dims":       VectorDims,
					"index":      true,
					"similarity": "cosine", // cosine similarity for normalized vectors
				},
				"product_id": map[string]any{"type": "keyword"},
				"branch_id":  map[string]any{"type": "keyword"},
			},
		},
		"settings": map[string]any{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
	}

	body, _ := json.Marshal(mapping)
	res, err := c.es.Indices.Create(IndexName, c.es.Indices.Create.WithContext(ctx), c.es.Indices.Create.WithBody(bytes.NewReader(body)))
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
		c.es.Index.WithOpType("create"),        // fail if doc already exists
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
