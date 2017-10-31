package qapi

import (
	"github.com/gorilla/websocket"
)

type websocketConnection struct {
	Conn *websocket.Conn
}

func (websocketConnection *websocketConnection) ReadQuotes() ([]Quote, error) {

	var response interface{}
	quotes := []Quote{}

	err := websocketConnection.Conn.ReadJSON(&response)
	if err !=nil {
		return nil, err
	}

	quotesTmp := response.(map[string]interface{})
	for k, v := range quotesTmp {
		switch k {
		case "quotes":
			for _, item := range v.([]interface{}) {
				quote := &Quote{}
				err := quote.FillStruct(item.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				quotes = append(quotes, *quote)
			}
		}
	}

	return quotes, nil
}
