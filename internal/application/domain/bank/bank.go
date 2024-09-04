package domain

import "time"

const (
	TransactionTypeIn  string = "IN"
	TransactionTypeOut string = "OUT"
)

type ExchangeRate struct {
	FromCurrency       string
	ToCurrency         string
	Rate               float64
	ValidFromTimestamp time.Time
	ValidToTimestamp   time.Time
}

type Transaction struct {
	Amount          float64
	TransactionType string
	Notes           string
}

type TransferTransaction struct {
	FromAccountNumber string
	ToAccountNumber   string
	Currency          string
	Amount            float64
}
