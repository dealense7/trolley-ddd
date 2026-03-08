package mysql

import (
	"github.com/dealense7/go-rates-ddd/internal/domain/country"
	"github.com/dealense7/go-rates-ddd/internal/domain/store"
)

type ReposContainer struct {
	CountryRepo country.Repository
	StoreRepo   store.Repository
}

// ProvideRepositories returns a container with all MySQL repositories.
func ProvideRepositories(db *DB) ReposContainer {
	return ReposContainer{
		CountryRepo: NewCountryRepo(db.DB),
		StoreRepo:   NewStoreRepo(db.DB),
	}
}
