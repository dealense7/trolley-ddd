package mysql

import (
	"github.com/dealense7/go-rates-ddd/internal/domain/country"
)

type ReposContainer struct {
	CountryRepo country.Repository
}

// ProvideRepositories returns a container with all MySQL repositories.
func ProvideRepositories(db *DB) ReposContainer {
	return ReposContainer{
		CountryRepo: NewCountryRepo(db.DB),
	}
}
