package util

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/TechTide8/CryptocurrencySDK/contract/erc20"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
)

type EthClient struct {
	*ethclient.Client
}

func NewEthClient(host string) *EthClient {
	rpcClient, err := rpc.Dial(host)
	if err != nil {
		panic(err)
	}
	ec := ethclient.NewClient(rpcClient)
	return &EthClient{
		Client: ec,
	}
}

func (*EthClient) CreateNewWallet() (priv, pub []byte, addr string, e error) {
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

func (client *EthClient) GetBalance(address common.Address) (decimal.Decimal, error) {
	ethAmount, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return decimal.Zero, err
	} else {
		balance := WeiToDecimal(ethAmount, ETHER)
		return balance, nil
	}
}

func (client *EthClient) GetTokenBalance(contractAddr common.Address, address common.Address, ethUnit EthUnit) (amount decimal.Decimal, err error) {
	token, err := erc20.NewErc20(contractAddr, client)
	if err != nil {
		return
	}
	balance, err := token.BalanceOf(nil, address)
	if err != nil {
		return
	}
	dec, err := token.Decimals(nil)

	if err != nil {
		amount = WeiToDecimal(balance, ethUnit)
		err = nil
	} else {
		amount = WeiToDecimal(balance, EthUnit(dec))
	}
	return
}

func (client *EthClient) Transfer(fromPri []byte, fromAddr, toAddr string, gasSpeedUp float64, amount decimal.Decimal) (txid string, maxGas decimal.Decimal, err error) {
	privateKey, err := crypto.ToECDSA(fromPri)
	if err != nil {
		return "", decimal.Zero, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", decimal.Zero, errors.New("PublicKey convert failed.")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", decimal.Zero, err
	}

	// gasPrice := new(big.Int).Div(ToWei(fee, ETHER), big.NewInt(int64(gasLimit)))
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", decimal.Zero, err
	}
	value := ToWei(amount, ETHER) // in wei (1 eth)

	toAddress := common.HexToAddress(toAddr)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		GasPrice: gasPrice,
		Value:    value,
		Data:     nil,
	})

	if err != nil {
		return "", decimal.Zero, err
	}

	gasLimit = uint64(float64(gasLimit) * gasSpeedUp)

	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", decimal.Zero, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", decimal.Zero, err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", decimal.Zero, err
	}

	txid = signedTx.Hash().Hex()
	maxGas = WeiToDecimal(gasPrice, ETHER).Mul(decimal.NewFromInt(int64(gasLimit)))

	return txid, maxGas, nil
}

func (client *EthClient) TransferToken(contractAddr string, fromPri []byte, fromAddr, toAddr string, tokenUnit string, gasSpeedUp float64, amount decimal.Decimal) (string, decimal.Decimal, error) {
	privateKey, err := crypto.ToECDSA(fromPri)
	if err != nil {
		return "", decimal.Zero, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", decimal.Zero, errors.New("PublicKey convert failed.")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	txid := ""

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", decimal.Zero, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", decimal.Zero, err
	}
	value := big.NewInt(0)
	toAddress := common.HexToAddress(toAddr)
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	amt := ToWei(amount, EthUnitFromString(tokenUnit))

	paddedAmount := common.LeftPadBytes(amt.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	tokenAddress := common.HexToAddress(contractAddr)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     fromAddress,
		To:       &tokenAddress,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	})
	if err != nil {
		return "", decimal.Zero, err
	}
	gasLimit = uint64(float64(gasLimit) * gasSpeedUp)

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", decimal.Zero, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", decimal.Zero, err
	}

	ctx := context.Background()
	err = client.SendTransaction(ctx, signedTx)

	if err != nil {
		return "", decimal.Zero, err
	}

	txid = signedTx.Hash().Hex()

	maxGas := WeiToDecimal(gasPrice, ETHER).Mul(decimal.NewFromInt(int64(gasLimit)))
	return txid, maxGas, nil
}
