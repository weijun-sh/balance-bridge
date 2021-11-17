package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	//"net/url"
	//"unsafe"
	"io/ioutil"
	"math"
	"math/big"
)

type ethConfig struct {
	Result string `bson: "result"`
	Error interface{} `bson: "error"`
}

func getBalance4ETH(url, address string) (string, bool) {
	//fmt.Printf("getBalance4ETH, url: %v, address: %v\n", url, address)
	data := make(map[string]interface{})
	data["method"] = "eth_getBalance"
	data["params"] = []string{address, "latest"}
	data["id"] = "1"
	data["jsonrpc"] = "2.0"
	bytesData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}
	basket := ethConfig{}
	for i := 0; i < 3; i++ {
		reader := bytes.NewReader(bytesData)
		resp, err := http.Post(url, "application/json", reader)
		if err != nil {
			fmt.Println(err.Error())
			return "", false
		}
		defer resp.Body.Close()

		//fmt.Printf("resp: %#v, resp.Body: %#v\n", resp, resp.Body)
		body, err := ioutil.ReadAll(resp.Body)
		//fmt.Printf("body: %v, string: %v\n", body, string(body))

		if err != nil {
			fmt.Println(err.Error())
			return "", false
		}
		err = json.Unmarshal(body, &basket)
		if err != nil {
			fmt.Println(err)
			return "", false
		}
		//fmt.Printf("%v basket.Result: %v, error: %v\n", i, basket.Result, basket.Error)
		if basket.Error != nil {
			//fmt.Printf("* Error *\n\n")
			basket.Error = nil
			continue
		} else {
			break
		}
	}
	b, g := getBalance4String(basket.Result, 18)
	return b, g
}

func getBalance4String(balance string, d int) (string, bool) {
	gas := false
	fbalance := new(big.Float)
	fbalance.SetString(balance)
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(d)))
	b, _ := ethValue.Float64()
	if b < minimumGas {
		gas = true
	}
	//f := fmt.Sprintf("%%18.%vf", d)
	//return fmt.Sprintf(f, b), gas
	return fmt.Sprintf("%0.4f", b), gas
}
