package wallet

import (
	"github.com/TechTide8/CryptocurrencySDK/model"
	"github.com/TechTide8/CryptocurrencySDK/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type ETH struct {
	client *util.EthClient
	conf   model.EvmConfig
}

func NewETH(host string, conf model.EvmConfig) *ETH {
	return &ETH{
		client: util.NewEthClient(host),
		conf:   conf,
	}
}

func (eth *ETH) NewAddress() (*model.Address, error) {
	pri, pub, address, err := eth.client.CreateNewWallet()
	if err != nil {
		return nil, err
	}

	result := &model.Address{
		Encryption: "Ethash",
		Address:    address,
		Pub:        pub,
		Pri:        pri,
	}
	return result, nil
}

func (eth *ETH) Transactions(trans []*model.ReqTransaction) []*model.ResTransaction {
	result := []*model.ResTransaction{}
	for _, record := range trans {
		txid, maxGas, err := eth.client.Transfer(
			record.SendAddress.Pri,
			record.SendAddress.Address,
			record.ToAddress,
			eth.conf.GasSpeedUp,
			record.Amount,
		)
		res := &model.ResTransaction{
			Uid:         record.Uid,
			SendAddress: record.SendAddress.Address,
			ToAddress:   record.ToAddress,
			Amount:      record.Amount,
		}
		if err == nil {
			res.Txid = txid
			res.EstimatedFee = maxGas
		} else {
			res.Error = err
		}
		result = append(result, res)
	}
	return result
}

func (eth *ETH) Balance(address string) (amount decimal.Decimal, err error) {
	return eth.client.GetBalance(common.HexToAddress(address))
}
