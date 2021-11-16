package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitClient(url *rpcConfig) *clientRPC {
	client := clientRPC{}
	client.C = make(map[string]*ethclient.Client, 0)
	//client.C["BTC"] = initClient(url.BTC)
	client.C["ETH"] = initClient(url.ETH)
	client.C["FSN"] = initClient(url.FSN)
	client.C["BSC"] = initClient(url.BSC)
	client.C["HT"] = initClient(url.HT)
	client.C["FTM"] = initClient(url.FTM)
	//client.C["LTC"] = initClient(url.LTC)
	//client.C["BLOCK"] = initClient(url.BLOCK)
	client.C["MATIC"] = initClient(url.MATIC)
	client.C["XDAI"] = initClient(url.XDAI)
	client.C["AVAX"] = initClient(url.AVAX)
	//client.C["HMY"] = initClient(url.HMY)
	//client.C["COLX"] = initClient(url.COLX)
	client.C["ARB"] = initClient(url.ARB)
	client.C["KCS"] = initClient(url.KCS)
	client.C["OKEX"] = initClient(url.OKEX)
	return &client
}

func initClient(gateway string) *ethclient.Client {
        ethcli, err := ethclient.Dial(gateway)
        if err != nil {
		fmt.Printf("ethclient.Dail failed, gateway: %v, err: %v\n", gateway, err)
		return nil
        }
	//fmt.Printf("ethclient.Dail gateway success, gateway: %v\n", gateway)
	return ethcli
}

type clientRPC struct {
	C map[string]*ethclient.Client
}

