package main

import (
	"fmt"
	"math/big"
	"os"
	"github.com/BurntSushi/toml"
)
var decimal map[string]*big.Int = make(map[string]*big.Int)
var emailTimeHour *[]uint64

func LoadDecimalConfig(filePath string) {
        if !FileExist(filePath) {
                fmt.Printf("LoadDecimalConfig error: config file '%v' not exist\n", filePath)
		return
        }

        config := &DecimalConfig{}
        if _, err := toml.DecodeFile(filePath, &config); err != nil {
                fmt.Printf("LoadDecimalConfig error (toml DecodeFile): %v\n", err)
		return
        }

	//fmt.Printf("LoadDecimalConfig, config.Contract: %v\n", config.Decimal)
	for _, t := range config.Decimal {
		chain := t.Chain
		token := t.Contract
		d := t.Decimal
		c := fmt.Sprintf("%v-%v", chain, token)
		decimal[c] = big.NewInt(d)
		//fmt.Printf("LoadDecimalConfig, decimal[c]: %v, d: %v\n", decimal[c], d)
	}
}

func LoadConfig(filePath string) *BridgeConfig {
        if !FileExist(filePath) {
                fmt.Printf("LoadConfig error: config file '%v' not exist\n", filePath)
		return nil
        }

        config := &BridgeConfig{}
        if _, err := toml.DecodeFile(filePath, &config); err != nil {
                fmt.Printf("LoadConfig error (toml DecodeFile): %v\n", err)
		return nil
        }

	emailTimeHour = config.Email.Time
        return config
}

func GetEmailTime() *[]uint64 {
	return emailTimeHour
}

type BridgeConfig struct {
	Email emailConfig
	Gas gasConfig
	Rpc rpcConfig
	Bridge []AddressConfig
}

type emailConfig struct {
	Time *[]uint64 // hour, 0-23
}

type gasConfig struct {
	MinimumGas float64
}

type AddressConfig struct {
	Msg string
	Important bool
	Address []string
	TokenAddress []string
}

type rpcConfig struct {
	BTC string
	ETH string
	FSN string
	BSC string
	HT string
	FTM string
	LTC string
	BLOCK string
	MATIC string
	XDAI string
	AVAX string
	HMY string
	COLX string
	ARB string
	KCS string
	OKEX string
	MOON string
	IOTEX string
	SHI string
	CELO string
	OETH string
	CRO string
	TLOS string
	TERRA string
	BOBA string
	FUSE string
}

type DecimalConfig struct {
	Decimal []ContractConfig
}

type ContractConfig struct {
	Chain string
	Contract string
	Decimal int64
}

// FileExist checks if a file exists at filePath.
func FileExist(filePath string) bool {
        _, err := os.Stat(filePath)
        if err != nil && os.IsNotExist(err) {
                return false
        }

        return true
}
