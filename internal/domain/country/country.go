package country

import "time"

type CurrencyCode string
type CurrencySymbol string

const (
	CurrencyCodeUSD   CurrencyCode   = "USD"
	CurrencyCodeEUR   CurrencyCode   = "EUR"
	CurrencyCodeGBP   CurrencyCode   = "GBP"
	CurrencyCodeGEL   CurrencyCode   = "GEL"
	CurrencySymbolUSD CurrencySymbol = "$"
	CurrencySymbolEUR CurrencySymbol = "€"
	CurrencySymbolGBP CurrencySymbol = "£"
	CurrencySymbolGEL CurrencySymbol = "₾"
)

type Country struct {
	ID int64 `db:"id"`

	Code           string         `db:"code"`
	Name           string         `db:"name"`
	NameLocal      string         `db:"name_local"`
	CurrencyCode   CurrencyCode   `db:"currency_code"`
	CurrencySymbol CurrencySymbol `db:"currency_symbol"`
	Timezone       string         `db:"timezone"`
	Active         bool           `db:"active"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewCountry(name, code, nameLocal, timezone string, currencyCode CurrencyCode, symbol CurrencySymbol) *Country {
	now := time.Now()
	return &Country{
		Code:           code,
		Name:           name,
		NameLocal:      nameLocal,
		CurrencyCode:   currencyCode,
		CurrencySymbol: symbol,
		Timezone:       timezone,
		Active:         true,
		CreatedAt:      now,
	}
}
