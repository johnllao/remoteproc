package models

import (
	"encoding/json"
)

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
