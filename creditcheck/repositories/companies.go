package repositories

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/johnllao/remoteproc/creditcheck/models"
)

var (
	ErrNotFound = errors.New("not found")
)

func loadCompanies(tx *bolt.Tx, companies *[]models.Company) error {
	var err error
	var b = tx.Bucket(BucketCompany)
	var c = b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var co *models.Company
		co, err = models.CompanyBytes(v).ToCompany()
		if err != nil {
			return err
		}
		*companies = append(*companies, *co)
	}
	return nil
}

func loadCompany(tx *bolt.Tx, symbol string, co *models.Company) error {
	var err error
	var b = tx.Bucket(BucketCompany)
	var d = b.Get([]byte(symbol))
	if d == nil {
		return ErrNotFound
	}
	var c *models.Company
	c, err = models.CompanyBytes(d).ToCompany()
	if err != nil {
		return err
	}
	*co = *c
	return nil
}

func addCompany(tx *bolt.Tx, c models.Company) error {
	var err error
	var b = tx.Bucket(BucketCompany)
	var v []byte
	v, err = c.ToJSON()
	if err != nil {
		return err
	}
	err = b.Put([]byte(c.Symbol), v)
	if err != nil {
		return err
	}
	return nil
}

func addCompanyToIndustryIndex(tx *bolt.Tx, industry, symbol, name string) error {
	var err error

	if industry == "" {
		return nil
	}
	if symbol == "" {
		return nil
	}

	var b = tx.Bucket(BucketIndustry)
	err = b.Put([]byte(industry), []byte(industry))
	if err != nil {
		return err
	}
	var k = append(BucketIndustry, []byte(":"+industry)...)
	b, err = b.CreateBucketIfNotExists(k)
	if err != nil {
		return err
	}
	err = b.Put([]byte(symbol), []byte(name))
	if err != nil {
		return err
	}
	return nil
}

func addCompanyToSectorIndex(tx *bolt.Tx, sector, symbol, name string) error {
	var err error

	if sector == "" {
		return nil
	}
	if symbol == "" {
		return nil
	}

	var b = tx.Bucket(BucketSector)
	err = b.Put([]byte(sector), []byte(sector))
	if err != nil {
		return err
	}
	var k = append(BucketSector, []byte(":"+sector)...)
	b, err = b.CreateBucketIfNotExists(k)
	if err != nil {
		return err
	}
	err = b.Put([]byte(symbol), []byte(name))
	if err != nil {
		return err
	}
	return nil
}
