package api

import (
	"fmt"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

type Sync struct {
	client  *marketdata.Client
	symbols []string
}

func NewSync(client *marketdata.Client) *Sync {
	return &Sync{
		client:  client,
		symbols: []string{"AAPL", "MSFT", "AMZN"},
	}
}

func (s *Sync) Start() bool {
	bars, err := s.client.GetLatestBars(s.symbols, marketdata.GetLatestBarRequest{})
	if err != nil {
		return false
	}

	for symbol, bar := range bars {
		fmt.Println(symbol)
		fmt.Println(bar.Close)
	}

	return true
}
