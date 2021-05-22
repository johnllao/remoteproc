package repositories

import (
	"fmt"

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
	var err error
	var tx *bolt.Tx
	tx, err = db.Begin(true)
	if err != nil {
		return nil
	}
	_, _ = tx.CreateBucketIfNotExists(BucketCompany)
	_, _ = tx.CreateBucketIfNotExists(BucketIndustry)
	_, _ = tx.CreateBucketIfNotExists(BucketSector)

	defer tx.Commit()

	var r = &Repository{
		DB: db,
	}

	return r
}

func (r *Repository) Companies() ([]models.Company, error) {
	var err error
	var companies = make([]models.Company, 0)
	err = r.DB.View(func(tx *bolt.Tx) error {
		err = loadCompanies(tx, &companies)
		if err != nil {
			return fmt.Errorf("WARN: Companies() failed to retrieve companies. err: %s", err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *Repository) FindCompany(symbol string) (*models.Company, error) {
	var err error
	var co models.Company
	err = r.DB.View(func(tx *bolt.Tx) error {
		err = loadCompany(tx, symbol, &co)
		if err == ErrNotFound {
			return err
		}
		return nil
	})
	if err == ErrNotFound {
		return nil, nil
	}
	return &co, err
}

func (r *Repository) SaveCompanies(companies []models.Company) error {
	return r.DB.Update(func(tx *bolt.Tx) error {
		var err error

		for _, c := range companies {
			// add the company record
			err = addCompany(tx, c)
			if err != nil {
				return fmt.Errorf("WARN: SaveCompanies() failed to save company. symbol: %s, err: %s", c.Symbol, err.Error())
			}

			// add the company to the industry index
			err = addCompanyToIndustryIndex(tx, c.Industry, c.Symbol, c.Name)
			if err != nil {
				return fmt.Errorf("WARN: SaveCompanies() failed to save company to industry index. symbol: %s, err: %s", c.Symbol, err.Error())
			}

			// add the company to the sector index
			err = addCompanyToSectorIndex(tx, c.Sector, c.Symbol, c.Name)
			if err != nil {
				return fmt.Errorf("WARN: SaveCompanies() failed to save company to sector index. symbol: %s, err: %s", c.Symbol, err.Error())
			}
		}

		return nil
	})
}
