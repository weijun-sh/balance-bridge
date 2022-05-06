package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"math"
	"math/big"
	"io"
	"os"
	"strings"
	"strconv"
	"time"
	com "github.com/weijun-sh/balance-bridge/common"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	//"github.com/syndtr/goleveldb/leveldb"
)

var (
	configFile *string
	errorFile *string
	decimalFile *string
	initDecimal *bool
	minimumGas float64
	insufficientGas bool
	sendEmail bool

	//db *leveldb.DB
	errf *os.File
	err error
	errRet string

	timeFrom = uint64(1629993600) // 20210827 00:00:00
)

func init() {
	configFile = flag.String("config", "", "config file")
	errorFile = flag.String("error", "error.txt", "log file")

	initDecimal = flag.Bool("initdecimal", false, "init decimal")
	decimalFile = flag.String("decimal", "", "decimal file")
	flag.Parse()

	errf, err = os.OpenFile(*errorFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		os.Exit(1)
	}
}

/*
/opt/balance-bridge/balance --config /opt/balance-bridge/config-balance.toml --decimal /opt/balance-bridge/config-decimal.toml > /opt/balance-bridge/bridge-withdraw.txt
*/

func main() {
	if *initDecimal {
		InitDecimal()
	}

	//db, err = leveldb.OpenFile("/opt/balance-bridge/db", nil)
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()

	config := LoadConfig(*configFile)
	minimumGas = config.Gas.MinimumGas
	LoadDecimalConfig(*decimalFile)
	//client := InitClient(&config.Rpc)

	emailH := GetEmailTime()
	length := len(config.Bridge)
	var i int = 0
	var j int = 0
	fmt.Printf("%v\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("ps. update every 30 minutes\n")
	fmt.Printf("ps. email time: %v o'clock everyday\n", *emailH)
	fmt.Printf("===============================================================\n")
	fmt.Printf("                   BALANCE  OF  FEE(BRIDGE)                    \n")
	fmt.Printf("===============================================================\n")
	for i = 0; i < length; i++ {
		if i > 0 {
			fmt.Printf("---------------------------------------------------------------\n")
		}
		b := config.Bridge[i]
		//fmt.Printf("msg: %v, addr: %v\n", b.Msg, b.Address)
		fmt.Printf("%v\n", b.Msg)
		errRet = fmt.Sprintf("%v\n", b.Msg)
		address := b.Address
		l := len(address)
		//fmt.Printf("l: %v\n", l)
		for j = 0; j < l; j++ {
			getBalance(config.Rpc, address[j], b.Important)
		}

		//tokenAddress := b.TokenAddress
		//lt := len(tokenAddress)
		//for j = 0; j < lt; j++ {
		//	getTokenBalance(client, tokenAddress[j], address[0])
		//}
		//fmt.Println()
	}
	fmt.Printf("===============================================================\n")
	fmt.Printf("* - important\n")
	fmt.Printf("%v\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("\n---- min GAS CONFIG ----\n")
	fmt.Printf("minGas: %v (if chain not set)\n", config.Gas.MinimumGas)
	fmt.Printf("chainMinGas: %v\n", config.Gas.ChainMinimumGas)
	defer errf.Close()
	//if sendEmail { // insufficientGas and balance changed
	//	os.Exit(1)
	//}
	if insufficientGas {
		nowHour := uint64(time.Now().Hour())
		for _, h := range *emailH {
			if nowHour == h {
				nowMinute := uint64(time.Now().Minute())
				if nowMinute < 30 {
					os.Exit(1)
				}
			}
		}
	}
}

func InitDecimal() {
	config := LoadConfig(*configFile)
	client := InitClient(&config.Rpc)
	length := len(config.Bridge)
	var i int = 0
	var j int = 0
	var e map[string]int = make(map[string]int)
	for i = 0; i < length; i++ {
		b := config.Bridge[i]
		slice := strings.Split(b.Address[0], ":")
		chain := slice[0]
		tokenAddress := b.TokenAddress
		lt := len(tokenAddress)
		for j = 0; j < lt; j++ {
			slicet := strings.Split(tokenAddress[j], ":")
			tokenAddr := slicet[1]
			chainU := strings.ToUpper(chain)
			s := fmt.Sprintf("%v-%v", chainU, tokenAddr)
			if e[s] == 1{
				continue
			}
			e[s] = 1
			decimal, errd := getErc20Decimal(client.C[chainU], tokenAddr)
			if  errd != nil {
				fmt.Printf("err, contract: %v, err: %v\n", tokenAddr, errd)
				continue
			}
			fmt.Printf("[[Decimal]]\nChain=\"%v\"\nContract=\"%v\"\nDecimal=%d\n\n", chain, tokenAddr, decimal)
		}
	}
}

func getTokenBalance(client *clientRPC, tokenAddress, address string) {
	slice := strings.Split(address, ":")
	chain := slice[0]
	addr := slice[1]
	slicet := strings.Split(tokenAddress, ":")
	name := slicet[0]
	tokenAddr := slicet[1]
	chainU := strings.ToUpper(chain)
	getTokenBalance4Chain(client, chainU, name, tokenAddr, addr)
}

func getDecimalBigFloat(d *big.Int) (*big.Float) {
	//fmt.Printf("getDecimalBigFloat, d: %v\n", d)
	f, _ := new(big.Float).SetInt(d).Float64()
	return big.NewFloat(math.Pow(10, f))
}

func getTokenBalance4Chain(client *clientRPC, chain, name, contract, addr string) {
	balance, errb := getErc20Balance(client.C[chain], contract, addr)
	s := fmt.Sprintf("%v-%v", chain, contract)
	//fmt.Printf("decimal[s]: %v\n", decimal[s])
	decimal := getDecimalBigFloat(decimal[s])
	if errb != nil {
		fmt.Printf("  %v  - %v\n", contract, name)
		return
	}
	b := new(big.Float).SetInt(balance)
	bd := new(big.Float).Quo(b, decimal)
	//b := convertBalance(balance)
	bdp := fmt.Sprintf("%0.4f", bd)
	fmt.Printf("  %v %15v %v\n", contract, bdp, name)
}

var erc20CodeParts = map[string][]byte{
        "name":         common.FromHex("0x06fdde03"),
        "symbol":       common.FromHex("0x95d89b41"),
        "decimal":      common.FromHex("0x313ce567"),
        "balanceOf":    common.FromHex("0x70a08231"),
}

// GetErc20Balance get erc20 decimal balacne of address
func getErc20Decimal(client *ethclient.Client, contract string) (*big.Int, error) {
        data := make([]byte, 4)
        copy(data[:4], erc20CodeParts["decimal"])
	to := common.HexToAddress(contract)
	msg := ethereum.CallMsg{
                To:   &to,
                Data: data,
        }
        result, err := client.CallContract(context.Background(), msg, nil)
        if err != nil {
                return nil, err
        }
	b := fmt.Sprintf("0x%v", hex.EncodeToString(result))
	return com.GetBigIntFromStr(b)
}

// GetErc20Balance get erc20 balacne of address
func getErc20Balance(client *ethclient.Client, contract, address string) (*big.Int, error) {
        data := make([]byte, 36)
        copy(data[:4], erc20CodeParts["balanceOf"])
        copy(data[4:], common.HexToAddress(address).Hash().Bytes())
	to := common.HexToAddress(contract)
	msg := ethereum.CallMsg{
                To:   &to,
                Data: data,
        }
        result, err := client.CallContract(context.Background(), msg, nil)
        if err != nil {
                return nil, err
        }
	b := fmt.Sprintf("0x%v", hex.EncodeToString(result))
        return com.GetBigIntFromStr(b)
}

func getBalance(rpc rpcConfig, address string, imp bool) {
	slice := strings.Split(address, ":")
	chain := slice[0]
	addr := slice[1]
	getBalance4Chain(rpc, chain, addr, imp)
	//fmt.Printf("chain: %v, addr: %v, rpc: %v\n", chain, addr, rpcURL)
}

func getBalance4Chain(rpc rpcConfig, chain, addr string, imp bool) {
	chainU := strings.ToUpper(chain)
	b := ""
	e := ""
	m := ""
	g := false
	switch chainU {
	case "BTC":
		b, m, g = getBalance4BTC(chainU, rpc.BTC, addr)
		e = fmt.Sprintf("  %v         %12v %v", addr, b, chainU)
	case "LTC":
		slice := strings.Split(addr, "-")
		addr1 := slice[0]
		addr2 := slice[1]
		b, m, g = getBalance4BTC(chainU, rpc.LTC, addr1)
		e = fmt.Sprintf("  %v         %12v %v", addr2, b, chainU)
	case "BLOCK":
		b, m, g = getBalance4BLOCK(chainU, rpc.BLOCK, addr)
		e = fmt.Sprintf("  %v         %12v %v", addr, b, chainU)
	case "COLX":
		slice := strings.Split(addr, "-")
		addr1 := slice[0]
		addr2 := slice[1]
		b, m, g  = getBalance4BTC(chainU, rpc.COLX, addr1)
		e = fmt.Sprintf("  %v         %12v %v", addr2, b, chainU)
	case "ETH":
		b, m, g = getBalance4ETH(chainU, rpc.ETH, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, chainU)
	case "FSN":
		b, m, g = getBalance4ETH(chainU, rpc.FSN, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, chainU)
	case "BSC":
		b, m, g = getBalance4ETH(chainU, rpc.BSC, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "BNB")
	case "FTM":
		b, m, g = getBalance4ETH(chainU, rpc.FTM, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, chainU)
	case "HT":
		b, m, g = getBalance4ETH(chainU, rpc.HT, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, chainU)
	case "MATIC":
		b, m, g = getBalance4ETH(chainU, rpc.MATIC, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, chainU)
	case "XDAI":
		b, m, g = getBalance4ETH(chainU, rpc.XDAI, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, chainU)
	case "AVAX":
		b, m, g = getBalance4ETH(chainU, rpc.AVAX, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, chainU)
	case "HMY":
		b, m, g = getBalance4ETH(chainU, rpc.HMY, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, chainU)
	case "ARB":
		b, m, g = getBalance4ETH(chainU, rpc.ARB, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "ETH(ARB)")
	case "KCS":
		b, m, g = getBalance4ETH(chainU, rpc.KCS, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, chainU)
	case "OKEX":
		b, m, g = getBalance4ETH(chainU, rpc.OKEX, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "OKT")
	case "MOON":
		b, m, g = getBalance4ETH(chainU, rpc.MOON, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "MOBR")
	case "IOTEX":
		b, m, g = getBalance4ETH(chainU, rpc.IOTEX, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "IOTX")
	case "SHI":
		b, m, g = getBalance4ETH(chainU, rpc.SHI, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "SDN")
	case "CELO":
		b, m, g = getBalance4ETH(chainU, rpc.CELO, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "CELO")
	case "OETH":
		b, m, g = getBalance4ETH(chainU, rpc.OETH, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "OETH")
	case "CRO":
		b, m, g = getBalance4ETH(chainU, rpc.CRO, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "CRO")
	case "TLOS":
		b, m, g = getBalance4ETH(chainU, rpc.TLOS, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "TLOS")
	case "TERRA":
		b, m, g = getBalance4ETH(chainU, rpc.TERRA, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "UST")
	case "BOBA":
		b, m, g = getBalance4ETH(chainU, rpc.BOBA, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "BOBA")
	case "FUSE":
		b, m, g = getBalance4ETH(chainU, rpc.FUSE, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "FUSE")
	case "SYS":
		b, m, g = getBalance4ETH(chainU, rpc.SYS, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "SYS")
	case "AURORA":
		b, m, g = getBalance4ETH(chainU, rpc.AURORA, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "ETH(AURORA)")
	case "METIS":
		b, m, g = getBalance4ETH(chainU, rpc.METIS, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "METIS")
	case "MOONBEAM":
		b, m, g = getBalance4ETH(chainU, rpc.MOONBEAM, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "GLMR")
	case "ASTAR":
		b, m, g = getBalance4ETH(chainU, rpc.ASTAR, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "ASTAR")
	case "ROSE":
		b, m, g = getBalance4ETH(chainU, rpc.ROSE, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "ROSE")
	case "VELAS":
		b, m, g = getBalance4ETH(chainU, rpc.VELAS, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "VLX")
	case "OASIS":
		b, m, g = getBalance4ETH(chainU, rpc.OASIS, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "ROSE")
	case "OPTIMISTIC":
		b, m, g = getBalance4ETH(chainU, rpc.OPTIMISTIC, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "OETH")
	case "CLV":
		b, m, g = getBalance4ETH(chainU, rpc.CLV, addr)
		e = fmt.Sprintf("  %v %12v %v", addr, b, "CLV")
	case "MIKO":
		b, m, g = getBalance4ETH(chainU, rpc.MIKO, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "MilkADA")
	case "NEBULAS":
		b, m, g = getBalance4ETH(chainU, rpc.NEBULAS, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "NEBULAS")
	case "REI":
		b, m, g = getBalance4ETH(chainU, rpc.REI, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "REI")
	case "CONFLUX":
		b, m, g = getBalance4ETH(chainU, rpc.CONFLUX, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "CFX")
	case "RSK":
		b, m, g = getBalance4ETH(chainU, rpc.RSK, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "RBTC")
	case "JEWEL":
		b, m, g = getBalance4ETH(chainU, rpc.JEWEL, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "JEWEL")
	case "BTTC":
		b, m, g = getBalance4ETH(chainU, rpc.BTTC, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "BTT")
	case "EVMOS":
		b, m, g = getBalance4ETH(chainU, rpc.EVMOS, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "EVMOS")
	case "RONIN":
		b, m, g = getBalance4ETH(chainU, rpc.RONIN, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "RON")
	case "CMP":
		b, m, g = getBalance4ETH(chainU, rpc.CMP, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "CMP")
	case "DOGE":
		b, m, g = getBalance4ETH(chainU, rpc.DOGE, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "DOGE")
	case "ETC":
		b, m, g = getBalance4ETH(chainU, rpc.ETC, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "ETC")
	case "XRP": // Different
		b, m, g = getBalance4XRP(chainU, rpc.XRP, addr)
		e = fmt.Sprintf("  %-42v %12v %v", addr, b, "XRP")
	default:
		return
	}
	if g {
		insufficientGas = true
		if imp {
			io.WriteString(errf, errRet)
			io.WriteString(errf, e)
			ee := fmt.Sprintf("      * ( < %v )\n", m)
			io.WriteString(errf, ee)
			e = fmt.Sprintf("%v     * ( < %v )", e, m)
		} else {
			e = fmt.Sprintf("%v     ( < %v )", e, m)
		}
		// leveldb
		//key := fmt.Sprintf("%v-%v", chain, addr)
		//balanceString := getKV(key)
		////fmt.Printf("b: %v, b-4db: %v\n", b, balanceString)
		//if balanceString != b {
		//	//fmt.Printf("!=\n")
		//	putKV(key, b)
		//	sendEmail = true
		//}
	}
	fmt.Println(e)
}

func Hex2Dec(val string) int {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return int(n)
}

