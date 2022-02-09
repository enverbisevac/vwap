//go:generate easyjson -all $GOFILE

package coinbase

// RequestType represents the type of request.
type RequestType string

const (
	// DefaultURL ...
	DefaultURL = "wss://ws-feed.exchange.coinbase.com"
)

const (
	// RequestTypeSubscribe ...
	RequestTypeSubscribe RequestType = "subscribe"
)

// ChannelType represents the type of channel on Coinbase.
type ChannelType string

// Channel ...
type Channel struct {
	Name ChannelType
}

// Request is a request to be sent to the Coinbase websocket.
type Request struct {
	Type       RequestType `json:"type"`
	ProductIDs []string    `json:"product_ids"`
	Channels   []Channel   `json:"channels"`
}

// Response is the response received after a request submission.
type Response struct {
	Type      string    `json:"type"`
	Channels  []Channel `json:"channels"`
	Message   string    `json:"message,omitempty"`
	Size      string    `json:"size"`
	Price     string    `json:"price"`
	ProductID string    `json:"product_id"`
}
