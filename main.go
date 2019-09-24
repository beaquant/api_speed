package main

import (
	"fmt"
	"github.com/beaquant/utils/json_file"
	"github.com/nntaoli-project/GoEx"
	"net/http"
)

type Config struct {
	Exchange string  `json:"exchange"`
	PubKey   string  `json:"pub_key"`
	SecKey   string  `json:"sec_key"`
	Pair     string  `json:"pair"`
	Amount   float64 `json:"amount"`
	Count    int     `json:"count"`
}

func main() {
	cfg := &Config{
		Exchange: "",
		PubKey:   "",
		SecKey:   "",
		Pair:     "",
		Count:    0,
	}
	json_file.Load("config.json", cfg)
	fmt.Println("cfg:", *cfg)
	pair := goex.NewCurrencyPair2(cfg.Pair)
	apiClient := NewApiClient(http.DefaultClient, cfg.PubKey, cfg.SecKey, cfg.Exchange, pair, cfg.Amount, cfg.Count)

	apiClient.GetTicker()
	apiClient.GetDepth()
	apiClient.GetAccount()
	apiClient.PutOrder()
	apiClient.CancelOrder()

	fmt.Println("get ticker:", apiClient.ApiState.GetTickerState.String())
	fmt.Println("get depth:", apiClient.ApiState.GetDepthState.String())
	fmt.Println("get account:", apiClient.ApiState.GetAccountState.String())
	fmt.Println("put order:", apiClient.ApiState.PutOrderState.String())
	fmt.Println("cancel order:", apiClient.ApiState.CancelOrderState.String())
}
