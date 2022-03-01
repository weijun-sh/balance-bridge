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

	//"github.com/davecgh/go-spew/spew"
)

type xrpConfig struct {
	Result resultConfig `json: "result"`
}

type resultConfig struct {
	Account_data accountDataConfig `json: "account_data"`
	Status string `json: "status"`
}

type accountDataConfig struct {
	Account string
	Balance string
}

type xrpParams struct {
	Account string `json: "account"`
	Strict bool `json: "strict"`
	Ledger_index string `json: "ledger_index"`
	Queue bool `json: "queue"`
}

func getBalance4XRP(url, address string) (string, bool) {
	//fmt.Printf("getBalance4XRP, url: %v, address: %v\n", url, address)
	data := "{\"method\":\"account_info\",\"params\":[{\"account\":\""+ address+"\",\"strict\":true,\"ledger_index\":\"current\",\"queue\":true}]}"
	basket := xrpConfig{}
	for i := 0; i < 1; i++ {
		reader := bytes.NewReader([]byte(data))
		resp, err := http.Post(url, "application/json", reader)
		if err != nil {
			fmt.Println(err.Error())
			return "", false
		}
		defer resp.Body.Close()

		//spew.Printf("resp.Body: %#v\n", resp.Body)
		body, err := ioutil.ReadAll(resp.Body)
		//spew.Printf("body: %#v, string: %v\n", body, string(body))

		if err != nil {
			fmt.Println(err.Error())
			return "", false
		}
		err = json.Unmarshal(body, &basket)
		if err != nil {
			fmt.Println(err)
			return "", false
		}
		//fmt.Printf("%v basket.Result: %v\n", i, basket.Result)
		if basket.Result.Status == "success" {
			break
		}
	}
	b, g := getBalance4String(basket.Result.Account_data.Balance, 6)
	return b, g
}

