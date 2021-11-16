package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type btcConfig struct {
	Address     string          `json: "address"`
	Chain_stats *btcStatsConfig `json: "chain_stats"`
	//mempool_stats interface{}
}

type btcStatsConfig struct {
	Funded_txo_count uint64 `json: "funded_txo_count"`
	Funded_txo_sum   uint64 `json: "funded_txo_sum"`
	Spent_txo_count  uint64 `json: "spent_txo_count"`
	Spent_txo_sum    uint64 `json: "spent_txo_sum"`
	Tx_count         uint64 `json: "tx_count"`
}

func getBalance4BTC(urlOrg, address string) (string, bool) {
	gas := false
	url := fmt.Sprintf("%v/address/%v", urlOrg, address)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))

	basket := btcConfig{}
	err = json.Unmarshal(body, &basket)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	funded := basket.Chain_stats.Funded_txo_sum
	spent := basket.Chain_stats.Spent_txo_sum
	balance := float64(funded-spent) / 100000000
	if balance < minimumGas {
		gas = true
	}
	return fmt.Sprintf("%0.4f", balance), gas
	//f := fmt.Sprintf("%%10.%vf", 8)
	//return fmt.Sprintf(f, balance), gas
}

type blockConfig struct {
	Utxos []utxoConfig `json: "utxos"`
}

type utxoConfig struct {
	Value float64 `json: "value"`
}

func getBalance4BLOCK(url, address string) (string, bool) {
	gas := false
	//fmt.Printf("getBalance4ETH, url: %v, address: %v\n", url, address)
	data := make(map[string]interface{})
	data["method"] = "getutxos"
	data["params"] = []string{"BLOCK", address}
	data["id"] = "1"
	data["version"] = "2.0"
	bytesData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}
	reader := bytes.NewReader(bytesData)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}
	basket := blockConfig{}
	err = json.Unmarshal(body, &basket)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	b := 0.0
	for _, v := range basket.Utxos {
		b += v.Value
	}
	if b != 0.0 && b < minimumGas {
		gas = true
	}
	return fmt.Sprintf("%0.4f", b), gas
}

