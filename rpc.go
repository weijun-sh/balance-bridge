package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitClient(gateway map[string]string) *clientRPC {
	client := clientRPC{}
	client.C = make(map[string]*ethclient.Client, 0)
	for id, url := range gateway {
		if doInitClient(id) {
			client.C[id] = initClient(url)
		}
	}
	return &client
}

func doInitClient(id string) bool {
	if id != "BTC" &&
	   id != "LTC" &&
	   id != "BLOCK" &&
	   id != "HMY" &&
	   id != "COLX" {
		return true
	}
	return false
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

