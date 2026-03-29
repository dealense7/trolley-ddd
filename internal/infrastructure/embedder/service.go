package embedder

import (
	"fmt"
)

type embedRequest struct {
	Name     string `json:"name,omitempty"`
	ImageB64 string `json:"image_b64,omitempty"`
}

type embedResponse struct {
	Vector []float64 `json:"vector"`
	Dims   int       `json:"dims"`
}

// EmbedFused generates a fused image+text vector (best accuracy).
func (c *Client) EmbedFused(imagePath, rawName string) ([]float64, error) {
	imgB64, err := imageToBase64(imagePath)
	if err != nil {
		return nil, fmt.Errorf("read image: %w", err)
	}
	return c.call(embedRequest{
		Name:     NormalizeName(rawName),
		ImageB64: imgB64,
	})
}

func (c *Client) embedImage(imagePath string) ([]float64, error) {
	imgB64, err := imageToBase64(imagePath)
	if err != nil {
		return nil, fmt.Errorf("read image: %w", err)
	}
	return c.call(embedRequest{ImageB64: imgB64})
}

func (c *Client) EmbedText(rawName string) ([]float64, error) {
	return c.call(embedRequest{Name: NormalizeName(rawName)})
}
