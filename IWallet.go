package CryptocurrencySDK

import (
	"github.com/TechTide8/CryptocurrencySDK/model"
	"github.com/TechTide8/CryptocurrencySDK/wallet"
	"github.com/shopspring/decimal"
)

var (
	_ IWallet = &wallet.ETH{}
)

type IWallet interface {
	// Create a wallet address.
	NewAddress() (*model.Address, error)
	// Transactions
	Transactions(trans []*model.ReqTransaction) []*model.ResTransaction
	// Get Wallet Balance on the Blockchain
	Balance(address string) (amount decimal.Decimal, err error)
}
