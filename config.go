package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	chainMinimumGas map[string]float64 = make(map[string]float64)
	decimal map[string]*big.Int = make(map[string]*big.Int)
	emailTimeHour *[]uint64
)

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

	initChainMinimumGas(config.Gas.ChainMinimumGas)
	emailTimeHour = config.Email.Time
        return config
}

func initChainMinimumGas(minGas []string) {
	for _, chainGas := range minGas {
		slice := strings.Split(chainGas, ":")
		if len(slice) != 2 {
			fmt.Printf("error: ChainMinimumGas '%v' length != 2\n", chainGas)
			continue
		}
		gasFloat, err := strconv.ParseFloat(slice[1], 64)
		if err != nil {
			fmt.Printf("error: ChainMinimumGas '%v' parse gas err: %v\n", chainGas, err)
			continue
		}
		chainU := strings.ToUpper(slice[0])
		chainMinimumGas[chainU] = gasFloat
	}
}

func GetChainMinimumGas(chain string) float64 {
	return chainMinimumGas[chain]
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
	ChainMinimumGas []string
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
	SYS string
	AURORA string
	METIS string
	MOONBEAM string
	ASTAR string
	ROSE string
	VELAS string
	OASIS string
	OPTIMISTIC string
	CLV string
	XRP string
	MIKO string
	NEBULAS string
	REI string
	CONFLUX string
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
