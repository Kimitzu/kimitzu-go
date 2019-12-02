package factory

import (
	"github.com/kimitzu/kimitzu-go/repo"
)

func NewSaleRecord() *repo.SaleRecord {
	contract := NewContract()
	return &repo.SaleRecord{
		Contract: contract,
		OrderID:  "anOrderIDforaSaleRecord",
	}
}
