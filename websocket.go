package qapi

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type WebsocketConnection struct {
	Conn *websocket.Conn
}

func (websocketConnection *WebsocketConnection) ReadQuotes() ([]Quote, error) {
	q := struct {
		Quotes []Quote `json:"quotes"`
	}{}

	err := websocketConnection.Conn.ReadJSON(&q)
	if err !=nil {
		return nil, errors.Wrap(err, "WebSocket connection failed to read json:\n")
	}

	return q.Quotes, nil
}
