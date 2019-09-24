package main

import (
	"fmt"
	"github.com/beaquant/utils"
	"github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/builder"
	"math"
	"net/http"
	"time"
)

const (
	Num = 30
)

type State struct {
	Max   time.Duration
	Min   time.Duration
	Ave   time.Duration
	Total time.Duration
	Count int64
}

type ApiState struct {
	PutOrderState    State
	CancelOrderState State
	GetAccountState  State
	GetTickerState   State
	GetDepthState    State
}

type ApiClient struct {
	API      goex.API
	Pair     goex.CurrencyPair
	ApiState ApiState
	Count    int
	orders   []goex.Order
}

func (s State) String() string {
	return fmt.Sprintf("Max:%s, Min:%s, Ave:%s, Total:%s, Count:%d", s.Max.String(), s.Min.String(), s.Ave.String(), s.Total.String(), s.Count)
}

func NewApiClient(httpClient *http.Client, pubKey, secKey, ex string, pair goex.CurrencyPair, count int) *ApiClient {
	return &ApiClient{
		API:    builder.NewCustomAPIBuilder(httpClient).APIKey(pubKey).APISecretkey(secKey).Build(ex),
		Pair:   pair,
		Count:  count,
		orders: make([]goex.Order, 0),
		ApiState: ApiState{
			PutOrderState: State{
				Max:   math.MinInt64,
				Min:   math.MaxInt64,
				Ave:   0,
				Total: 0,
				Count: 0,
			},
			CancelOrderState: State{
				Max:   math.MinInt64,
				Min:   math.MaxInt64,
				Ave:   0,
				Total: 0,
				Count: 0,
			},
			GetAccountState: State{
				Max:   math.MinInt64,
				Min:   math.MaxInt64,
				Ave:   0,
				Total: 0,
				Count: 0,
			},
			GetTickerState: State{
				Max:   math.MinInt64,
				Min:   math.MaxInt64,
				Ave:   0,
				Total: 0,
				Count: 0,
			},
			GetDepthState: State{
				Max:   math.MinInt64,
				Min:   math.MaxInt64,
				Ave:   0,
				Total: 0,
				Count: 0,
			},
		},
	}
}

func (a *ApiClient) updateTime(state *State, cost time.Duration) {
	if state.Max < cost {
		state.Max = cost
	}
	if state.Min > cost {
		state.Min = cost
	}
	state.Count++
	state.Total += cost
	state.Ave = state.Total / time.Duration(state.Count)
}

func (a *ApiClient) PutOrder() {
	fmt.Println("start put order test...")
	ticker, err := a.API.GetTicker(a.Pair)
	if err != nil {
		return
	}

	price := ticker.Buy * (1 - 0.095)
	for i := 0; i < a.Count; i++ {
		start := time.Now()
		ord, err := a.API.LimitBuy("0.1", utils.Float64RoundString(price), a.Pair)
		if err != nil {
			fmt.Println("PutOrder err:", err)
			continue
		}
		cost := time.Since(start)
		a.updateTime(&a.ApiState.PutOrderState, cost)
		a.orders = append(a.orders, *ord)
		time.Sleep(time.Second)
	}
	fmt.Println("finish put order test...")
}

func (a *ApiClient) CancelOrder() {
	fmt.Println("start cancel order test...")
	for _, v := range a.orders {
		start := time.Now()
		_, err := a.API.CancelOrder(v.OrderID2, a.Pair)
		if err == nil {
			cost := time.Since(start)
			a.updateTime(&a.ApiState.CancelOrderState, cost)
		} else {
			fmt.Println("CancelOrder err:", err)
		}
		time.Sleep(time.Second)
	}
	fmt.Println("finish cancel order test...")
}

func (a *ApiClient) GetAccount() {
	fmt.Println("start get account test...")
	for i := 0; i < a.Count; i++ {
		start := time.Now()
		_, err := a.API.GetAccount()
		if err != nil {
			fmt.Println("GetAccount err:", err)
			continue
		}
		cost := time.Since(start)
		a.updateTime(&a.ApiState.GetAccountState, cost)
		time.Sleep(time.Second)
	}
	fmt.Println("finish get account test...")
}

func (a *ApiClient) GetTicker() {
	fmt.Println("start get ticker test...")
	for i := 0; i < a.Count; i++ {
		start := time.Now()
		_, err := a.API.GetTicker(a.Pair)
		if err != nil {
			fmt.Println("GetTicker err:", err)
			continue
		}
		cost := time.Since(start)
		a.updateTime(&a.ApiState.GetTickerState, cost)
		time.Sleep(time.Second)
	}
	fmt.Println("finish get ticker test...")
}

func (a *ApiClient) GetDepth() {
	fmt.Println("start get depth test...")
	for i := 0; i < a.Count; i++ {
		start := time.Now()
		_, err := a.API.GetDepth(5, a.Pair)
		if err != nil {
			fmt.Println("GetDepth err:", err)
			continue
		}
		cost := time.Since(start)
		a.updateTime(&a.ApiState.GetDepthState, cost)
		time.Sleep(time.Second)
	}
	fmt.Println("finish get depth test...")
}
