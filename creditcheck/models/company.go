package models

import (
	"encoding/json"
)

type CompanyBytes []byte

type Company struct {
	Symbol        string
	Name          string
	LastSale      string
	NetChange     string
	PercentChange string
	MarketCap     string
	Country       string
	IPOYear       int
	Volume        int
	Sector        string
	Industry      string
}

func (c Company) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

func (b CompanyBytes) ToCompany() (*Company, error) {
	var err error
	var c = new(Company)
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
