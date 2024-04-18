package wallet

import (
	"crypto/ecdsa"
	"errors"

	"github.com/TechTide8/CryptocurrencySDK/model"
	"github.com/ethereum/go-ethereum/crypto"
)

type ETH struct {
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
