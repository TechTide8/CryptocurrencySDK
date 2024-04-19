package util

import (
	"math/big"
	"reflect"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type EthUnit int

const (
	WEI    EthUnit = 1
	KWEI           = 3
	WWEI           = 4
	MWEI           = 6
	LWEI           = 8
	GWEI           = 9
	SZABO          = 12
	FINNEY         = 15
	ETHER          = 18
	KETHER         = 21
	METHER         = 24
	GETHER         = 27
)

func EthUnitFromString(unit string) EthUnit {
	switch unit {
	case "WEI":
		return WEI
	case "KWEI":
		return KWEI
	case "WWEI":
		return WWEI
	case "MWEI":
		return MWEI
	case "LWEI":
		return LWEI
	case "GWEI":
		return GWEI
	case "SZABO":
		return SZABO
	case "FINNEY":
		return FINNEY
	case "ETHER":
		return ETHER
	case "KETHER":
		return KETHER
	case "METHER":
		return METHER
	case "GETHER":
		return GETHER
	default:
		return ETHER
	}
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// IsZeroAddress validate if it's a 0 address
func IsZeroAddress(iaddress interface{}) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

// ToDecimal wei to decimals
func WeiToDecimal(ivalue *big.Int, unit EthUnit) decimal.Decimal {

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromInt32(int32(unit)))
	num, _ := decimal.NewFromString(ivalue.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func ToWei(iamount interface{}, unit EthUnit) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromInt32(int32(unit)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// CalcGasCost calculate gas cost given gas limit (units) and gas price (wei)
func CalcGasCost(gasLimit uint64, gasPrice *big.Int) *big.Int {
	gasLimitBig := big.NewInt(int64(gasLimit))
	return gasLimitBig.Mul(gasLimitBig, gasPrice)
}
