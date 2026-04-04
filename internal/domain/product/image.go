package product

import "time"

type Image struct {
	ID int64

	ProductId int64 `db:"product_id"`

	Name      string `db:"name"`
	Size      int64  `db:"size"`
	Extension string `db:"extension"`
	Folder    string `db:"folder"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
