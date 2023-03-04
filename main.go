package main

import (
	"log"
	"os"

	"github.com/MinimalistMarkets/mrkts.fyi/api"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/joho/godotenv"
)

func main() {
	log.Print("starting server...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file")
	}

	client := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    os.Getenv("ALPACA_API_KEY"),
		APISecret: os.Getenv("ALPACA_API_SECRET"),
		BaseURL:   os.Getenv("ALPACA_DATA_BASE_URL"),
	})

	// Alpaca market data client has no way to verify that the API key and secret are valid.
	// This is a workaround to verify that the client is configured correctly.
	if _, err = client.GetLatestBar("AAPL", marketdata.GetLatestBarRequest{}); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(":"+port, client)
	log.Printf("listening on port %s", port)
	if err = server.Start(); err != nil {
		log.Fatal(err)
	}
}
