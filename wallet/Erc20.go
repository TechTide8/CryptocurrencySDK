package wallet

import (
	"github.com/TechTide8/CryptocurrencySDK/model"
	"github.com/TechTide8/CryptocurrencySDK/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type Erc20 struct {
	client       *util.EthClient
	conf         model.EvmConfig
	contractConf model.EvmContract
}

func NewErc20(host string, conf model.EvmConfig, contractConf model.EvmContract) *Erc20 {
	return &Erc20{
		client:       util.NewEthClient(host),
		conf:         conf,
		contractConf: contractConf,
	}
}

func (eth *Erc20) NewAddress() (*model.Address, error) {
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

func (eth *Erc20) Transactions(trans []*model.ReqTransaction) []*model.ResTransaction {
	result := []*model.ResTransaction{}
	for _, record := range trans {
		txid, maxGas, err := eth.client.TransferToken(
			eth.contractConf.ContractAddress,
			record.SendAddress.Pri,
			record.SendAddress.Address,
			record.ToAddress,
			eth.conf.Symbol,
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

func (eth *Erc20) Balance(address string) (amount decimal.Decimal, err error) {
	return eth.client.GetTokenBalance(
		common.HexToAddress(eth.contractConf.ContractAddress),
		common.HexToAddress(address),
		util.EthUnitFromString(eth.contractConf.Unit))
}
