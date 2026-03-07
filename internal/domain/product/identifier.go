package product

import "time"

type IdentifierType string

const (
	IdentifierTypeBarcode          IdentifierType = "barcode"
	IdentifierTypeGTIN             IdentifierType = "gtin"
	IdentifierTypeEAN13            IdentifierType = "ean13"
	IdentifierTypeUPC              IdentifierType = "upc"
	IdentifierTypeSKU              IdentifierType = "sku"
	IdentifierTypeStoreID          IdentifierType = "store_id"
	IdentifierTypeManufacturerCode IdentifierType = "manufacturer_code"
)

type Identifier struct {
	ID        int64
	ProductId int64 `db:"product_id"`

	Type  IdentifierType `db:"type"`
	Value string         `db:"value"`

	Confidence string `db:"confidence"`

	CreatedAt time.Time `db:"created_at"`
}
