package model

import "github.com/shopspring/decimal"

type ResTransaction struct {
	Uid          uint            // Unique custom ID
	Txid         string          // Transaction ID on the blockchain
	SendAddress  string          // From wallet address
	ToAddress    string          // Recipient wallet address
	Amount       decimal.Decimal // Received amount
	EstimatedFee decimal.Decimal // Estimated Transaction Fee
	Error        error
}
