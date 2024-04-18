package wallet

import (
	"crypto/ecdsa"
	"errors"

	"github.com/TechTide8/CryptocurrencySDK/model"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shopspring/decimal"
)

type ETH struct {
	client *ethclient.Client
}

func NewETH(host string) *ETH {
	rpcClient, err := rpc.Dial(host)
	if err != nil {
		panic(err)
	}
	ec := ethclient.NewClient(rpcClient)
	return &ETH{
		client: ec,
	}
}

func (eth *ETH) NewAddress() (*model.Address, error) {
	pri, pub, address, err := eth.createNewWallet()
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
	return nil
}

func (eth *ETH) Balance(address string) (amount decimal.Decimal, err error) {
	return decimal.Zero, nil
}

func (*ETH) createNewWallet() (priv, pub []byte, addr string, e error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil, "", err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")

	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return privateKeyBytes, publicKeyBytes, address, nil
}
