package embedder

import (
	"net/http"
	"time"

	"github.com/dealense7/go-rates-ddd/internal/common/cfg"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(cfg *cfg.Config) *Client {
	return &Client{
		baseURL:    cfg.EmbedderURL,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}
