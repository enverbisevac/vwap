//go:generate easyjson -all $GOFILE

package websocket

import (
	"context"
	"encoding/json"

	ws "golang.org/x/net/websocket"
)

func init() {
	ws.JSON = ws.Codec{
		Marshal:   jsonMarshal,
		Unmarshal: jsonUnmarshal,
	}
}

// Response is the response received after a request submission.
type Response struct {
	Type      string `json:"type"`
	Size      string `json:"size"`
	Price     string `json:"price"`
	ProductID string `json:"product_id"`
}

type Subscriber interface {
	Subscribe(ctx context.Context, tradingPairs []string, receiver chan Response) error
}

func jsonMarshal(v interface{}) (msg []byte, payloadType byte, err error) {
	msg, err = v.(json.Marshaler).MarshalJSON()
	return msg, ws.TextFrame, err
}

func jsonUnmarshal(msg []byte, payloadType byte, v interface{}) (err error) {
	return v.(json.Unmarshaler).UnmarshalJSON(msg)
}
