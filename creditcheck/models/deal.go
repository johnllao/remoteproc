package models

import (
	"encoding/json"
)

type DealBytes []byte

type Deal struct {
	TxnID  string  `json:"txnid"`
	Symbol string  `json:"symbol"`
	Amount float64 `json:"amt"`
}

func (d Deal) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

func (d DealBytes) ToDeal() (*Deal, error) {
	var err error
	var deal = new(Deal)
	err = json.Unmarshal(d, deal)
	if err != nil {
		return nil, err
	}
	return deal, nil
}
