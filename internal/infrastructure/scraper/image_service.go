package scraper

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chai2010/webp"
	"github.com/dealense7/go-rates-ddd/internal/domain/product"
	"github.com/google/uuid"
)

func (s *ParserService) downloadImage(productId int64, url string) (*product.Image, error) {
	// 1. Download image
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// 2. Read image into buffer (so we can decode multiple times)
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(buf)

	// 3. Try normal image decode (JPEG/PNG/GIF)
	img, _, err := image.Decode(reader)
	if err != nil {
		// reset reader and try WebP
		reader.Seek(0, io.SeekStart)
		img, err = webp.Decode(reader)
		if err != nil {
			return nil, fmt.Errorf("cannot decode image: %w", err)
		}
	}

	// 4. Prepare folder and UUID filename
	folder := "static/images/products"
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return nil, err
	}
	id := uuid.New().String()
	filePath := filepath.Join(folder, id+".webp")

	// 5. Save as WebP
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := webp.Encode(f, img, &webp.Options{Lossless: true}); err != nil {
		return nil, err
	}

	// 6. File info
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &product.Image{
		ProductId: productId,
		Name:      id,
		Size:      info.Size(),
		Extension: "webp",
		Folder:    folder,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
