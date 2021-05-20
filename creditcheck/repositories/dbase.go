package repositories

import (
	"github.com/boltdb/bolt"

	"github.com/johnllao/remoteproc/creditcheck/models"
)

var (
	BucketCompany,
	BucketIndustry,
	BucketSector []byte
)

type Repository struct {
	DB *bolt.DB
}

func init() {
	BucketCompany = []byte("COMPANY")
	BucketIndustry = []byte("INDUSTRY")
	BucketSector = []byte("SECTOR")
}

func NewRepository(db *bolt.DB) *Repository {
	var r = &Repository{
		DB: db,
	}

	return r
}

func (r *Repository) SaveCompanies(companies []*models.Company) error {
	return r.DB.Update(func(tx *bolt.Tx) error {
		var err error

		for _, c := range companies {
			// add the company record
			err = addCompany(tx, c)
			if err != nil {
				return err
			}

			// add the company to the industry index
			err = addCompanyToIndustryIndex(tx, c.Industry, c.Symbol, c.Name)
			if err != nil {
				return err
			}

			// add the company to the sector index
			err = addCompanyToSectorIndex(tx, c.Sector, c.Symbol, c.Name)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
