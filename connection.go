package qapi

import (
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type WebsocketConnection struct {
	Conn *websocket.Conn
}

func (websocketConnection *WebsocketConnection) ReadQuotes() ([]Quote, error) {
	var response interface{}
	quotes := []Quote{}

	err := websocketConnection.Conn.ReadJSON(&response)
	if err !=nil {
		return nil, errors.Wrap(err, "WebSocket connection failed to read json:\n")
	}

	quotesTmp := response.(map[string]interface{})
	for k, v := range quotesTmp {
		if k == "quotes" {
			for _, item := range v.([]interface{}) {
				quote := &Quote{}
				// Fill struct quote with item's data
				err := mapstructure.Decode(item.(map[string]interface{}), &quote)
				if err != nil {
					return nil, errors.Wrap(err, "Unable to decode data into Quote struct:\n")
				}
				quotes = append(quotes, *quote)
			}
		}
	}

	return quotes, nil
}
