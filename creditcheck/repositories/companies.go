package repositories

import (
	"encoding/binary"
	"errors"
	"math"

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

func loadCompanyLimitAntUtil(tx *bolt.Tx, symbol string, limit, util *float64) error {
	var limb = tx.Bucket(BucketLimits).Get([]byte(symbol))
	var utilb = tx.Bucket(BucketUtilizations).Get([]byte(symbol))

	*limit = 0
	if limb != nil {
		*limit = math.Float64frombits(binary.LittleEndian.Uint64(limb))
	}
	*util = 0
	if utilb != nil {
		*util = math.Float64frombits(binary.LittleEndian.Uint64(utilb))
	}
	return nil
}

func updateCompanyLimit(tx *bolt.Tx, symbol string, limit float64) error {
	var err error
	var b = tx.Bucket(BucketLimits)
	var l = make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(l, math.Float64bits(limit))
	err = b.Put([]byte(symbol), l)
	if err != nil {
		return err
	}
	return nil
}

func updateCompanyUtilization(tx *bolt.Tx, symbol string, util float64) error {
	var err error
	var b = tx.Bucket(BucketUtilizations)
	var u = make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(u, math.Float64bits(util))
	err = b.Put([]byte(symbol), u)
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
