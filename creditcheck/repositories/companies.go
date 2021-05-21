package repositories

import (
	"github.com/boltdb/bolt"

	"github.com/johnllao/remoteproc/creditcheck/models"
)

func addCompany(tx *bolt.Tx, c *models.Company) error {
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
