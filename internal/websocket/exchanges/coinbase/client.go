package coinbase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/enverbisevac/vwap/internal/websocket"

	ws "golang.org/x/net/websocket"
)

type Client struct {
	conn *ws.Conn
}

// NewClient returns a new websocket client.
func NewClient(url string) (Client, error) {
	conn, err := ws.Dial(url, "", "http://localhost/")
	if err != nil {
		return Client{}, err
	}

	log.Printf("websocket connected to: %s", url)

	return Client{
		conn: conn,
	}, nil
}

// Subscribe subscribes to the websocket.
func (c Client) Subscribe(ctx context.Context, tradingPairs []string, receiver chan websocket.Response) error {
	if len(tradingPairs) == 0 {
		return errors.New("tradingPairs must be provided")
	}

	subscription := Request{
		Type:       RequestTypeSubscribe,
		ProductIDs: tradingPairs,
		Channels: []Channel{
			{Name: "matches"},
		},
	}

	payload, err := json.Marshal(subscription)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %w", err)
	}

	err = ws.Message.Send(c.conn, payload)
	if err != nil {
		return fmt.Errorf("failed to send subscription: %w", err)
	}

	subscriptionResponse := Response{}

	err = ws.JSON.Receive(c.conn, &subscriptionResponse)
	if err != nil {
		return fmt.Errorf("failed to receive subscription response: %w", err)
	}

	if subscriptionResponse.Type == "error" {
		return fmt.Errorf("failed to subscribe: %s", subscriptionResponse.Message)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := c.conn.Close()
				if err != nil {
					log.Printf("failed closing ws connection: %s", err)
				}
				log.Println("exiting from go routine")
				close(receiver)
				return
			default:
				message := websocket.Response{}

				err := ws.JSON.Receive(c.conn, &message)
				if err != nil {
					log.Printf("failed receiving message: %s", err)
					break
				}

				receiver <- message
			}
		}
	}()

	return nil
}
