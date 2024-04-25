package model

type EvmConfig struct {
	ChainNetwork string
	Symbol       string
	Host         string
	GasSpeedUp   float64
	Contracts    []*EvmContract
}

type EvmContract struct {
	Symbol          string
	ContractAddress string
	Unit            string
}
