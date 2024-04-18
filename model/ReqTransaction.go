package model

import "github.com/shopspring/decimal"

type ReqTransaction struct {
	Uid         uint            // Unique custom ID
	SendAddress Address         // From wallet information
	ToAddress   string          // Recipient wallet address
	Txid        string          // Transaction ID on the blockchain
	CoinKey     string          // Cryptocurrency symbol
	ChainKey    string          // The blockchain name, eg. ETH、BTC、Tron、Binance
	Amount      decimal.Decimal // Received amount
	Fee         decimal.Decimal // Transaction fee on the blockchain, denominated in the base currency of the main chain.
}
