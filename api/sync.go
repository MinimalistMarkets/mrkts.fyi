package api

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/MinimalistMarkets/mrkts.fyi/storage"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

type PathDataMap map[string][]string

type Sync struct {
	client  *marketdata.Client
	storage *storage.CloudStorage
	symbols PathDataMap
}

type BasicTickerData struct {
	Symbol string
	Price  float64
}

func NewSync(client *marketdata.Client, storage *storage.CloudStorage) *Sync {
	pathDataMap := PathDataMap{
		"index.html": []string{"AAPL", "MSFT", "AMZN", "GOOG", "FB", "TSLA", "NFLX", "NVDA", "PYPL", "ADBE"},
	}
	return &Sync{
		client:  client,
		storage: storage,
		symbols: pathDataMap,
	}
}

func (s *Sync) Start() bool {
	ok := s.SyncPage("index.html")
	return ok
}

func (s *Sync) SyncPage(path string) bool {
	data := s.GetTickerData(path)
	if data == nil {
		log.Printf("failed to get ticker data: %s", path)
		return false
	}

	filename, err := s.SaveDataToHTML(data, path)
	if err != nil {
		log.Printf("failed to save data to html: %s - %v", path, err)
		return false
	}
	defer os.Remove(filename)

	if err := s.storage.Upload(filename, path); err != nil {
		return false
	}

	return true
}

func (s *Sync) GetTickerData(path string) []BasicTickerData {
	bars, err := s.client.GetLatestBars(s.symbols[path], marketdata.GetLatestBarRequest{})
	if err != nil {
		return nil
	}

	var data []BasicTickerData
	for symbol, bar := range bars {
		data = append(data, BasicTickerData{
			Symbol: symbol,
			Price:  bar.Close,
		})
	}

	return data
}

func (s *Sync) SaveDataToHTML(data []BasicTickerData, path string) (filename string, err error) {
	t, _ := template.ParseFiles("templates/base.html", "templates/"+path)

	file, err := ioutil.TempFile("./", strings.Split(path, ".")[0]+"-*.html")
	if err != nil {
		return "", err
	}

	err = t.Execute(file, data)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
