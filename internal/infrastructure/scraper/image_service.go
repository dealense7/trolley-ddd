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
	"golang.org/x/image/draw"
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

	// 4. Resize
	img = resizeAndCenterCrop(img, 224, 224)

	// 5. Prepare folder and UUID filename
	folder := "static/images/products"
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return nil, err
	}
	id := uuid.New().String()
	filePath := filepath.Join(folder, id+".webp")

	// 6. Save as WebP
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := webp.Encode(f, img, &webp.Options{Lossless: false, Quality: 85}); err != nil {
		return nil, err
	}

	// 7. File info
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
	}, nil
}

func resizeAndCenterCrop(src image.Image, targetW, targetH int) image.Image {
	bounds := src.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()

	// Scale so the short side fills the target
	scaleW := float64(targetW) / float64(srcW)
	scaleH := float64(targetH) / float64(srcH)
	scale := scaleW
	if scaleH > scaleW {
		scale = scaleH
	}

	scaledW := int(float64(srcW) * scale)
	scaledH := int(float64(srcH) * scale)

	// 1. Resize
	scaled := image.NewRGBA(image.Rect(0, 0, scaledW, scaledH))
	draw.BiLinear.Scale(scaled, scaled.Bounds(), src, bounds, draw.Over, nil)

	// 2. Center crop
	x0 := (scaledW - targetW) / 2
	y0 := (scaledH - targetH) / 2
	cropped := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	draw.Draw(cropped, cropped.Bounds(), scaled, image.Point{x0, y0}, draw.Src)

	return cropped
}
