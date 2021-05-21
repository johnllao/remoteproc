package models

import (
	"encoding/json"
)

type CompanyBytes []byte

type Company struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	LastSale      string  `json:"lastsale"`
	NetChange     float64 `json:"netchg"`
	PercentChange string  `json:"pctchg"`
	MarketCap     float64 `json:"mktcap"`
	Country       string  `json:"country"`
	IPOYear       int     `json:"ipoyr"`
	Volume        int     `json:"vol"`
	Sector        string  `json:"sec"`
	Industry      string  `json:"ind"`
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
