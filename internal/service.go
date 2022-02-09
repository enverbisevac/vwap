package internal

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/enverbisevac/vwap/internal/websocket"

	"github.com/enverbisevac/vwap/internal/vwap"
)

// Service is our main service.
type Service struct {
	wsClient     websocket.Subscriber
	tradingPairs []string
	list         map[string]*vwap.Buffer
}

// NewService returns a new service.
func NewService(ws websocket.Subscriber, tradingPairs []string, defaultSize uint) Service {
	m := make(map[string]*vwap.Buffer, len(tradingPairs))
	for _, val := range tradingPairs {
		m[val] = vwap.NewBuffer(defaultSize)
	}
	return Service{
		wsClient:     ws,
		tradingPairs: tradingPairs,
		list:         m,
	}
}

// Run service
func (s Service) Run(ctx context.Context) error {
	receiver := make(chan websocket.Response)
	defer log.Println("No more data service stopped")

	err := s.wsClient.Subscribe(ctx, s.tradingPairs, receiver)
	if err != nil {
		return fmt.Errorf("service subscription err: %w", err)
	}

	for data := range receiver {
		if data.Price == "" {
			continue
		}

		decimalPrice, err := strconv.ParseFloat(data.Price, 64)
		if err != nil {
			return fmt.Errorf("decimalPrice %s: %w", data.Price, err)
		}

		decimalSize, err := strconv.ParseFloat(data.Size, 64)
		if err != nil {
			return fmt.Errorf("decimalSize %s: %w", data.Size, err)
		}

		s.list[data.ProductID].Push(vwap.DataPoint{
			Price:  decimalPrice,
			Volume: decimalSize,
		})

		// Print to sdout.
		fmt.Println(data.ProductID, s.list[data.ProductID].VWAP)
	}

	return nil
}
