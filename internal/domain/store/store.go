package store

import "time"

type Store struct {
	ID int64 `db:"id"`

	Name         string `db:"name"`
	Slug         string `db:"slug"`
	LogoURL      string `db:"logo_url"`
	PrimaryColor string `db:"primary_color"`
	Active       bool   `db:"active"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Branches *[]Branch `db:"-"`
}

func NewStore(name, slug, logoUrl, primaryColor string, branches *[]Branch) *Store {
	now := time.Now()
	return &Store{
		Name:         name,
		Slug:         slug,
		LogoURL:      logoUrl,
		PrimaryColor: primaryColor,
		Active:       true,
		CreatedAt:    now,
		Branches:     branches,
	}
}
