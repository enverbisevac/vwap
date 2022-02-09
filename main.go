package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/enverbisevac/vwap/internal"
	"github.com/enverbisevac/vwap/internal/websocket/exchanges/coinbase"
)

const (
	defaultTradingPairs = "BTC-USD,ETH-USD,ETH-BTC"
	defaultWindowSize   = 200
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ws, err := coinbase.NewClient(coinbase.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	service := internal.NewService(&ws, strings.Split(defaultTradingPairs, ","), defaultWindowSize)

	err = service.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
