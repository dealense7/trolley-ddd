package product

import (
	"fmt"
	"time"
)

type Image struct {
	ID int64

	ProductId int64 `db:"product_id"`

	Name          string `db:"name"`
	Size          int64  `db:"size"`
	Extension     string `db:"extension"`
	Folder        string `db:"folder"`
	HasEmbeddings bool   `db:"has_embeddings"`

	CreatedAt time.Time `db:"created_at"`
}

func (i *Image) ImageURL() string {
	return fmt.Sprintf("%s/%s.%s", i.Folder, i.Name, i.Extension)
}
