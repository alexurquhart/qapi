package qapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type QuestradeError struct {
	Code       int `json:"code",string`
	StatusCode int
	Message    string `json:"message"`
	Endpoint   string
	OrderId    int     `json:"orderId,omitempty"`
	Orders     []Order `json:"orders,omitempty"`
}

func newQuestradeError(res *http.Response, body []byte) QuestradeError {
	// Unmarshall the error text
	var e QuestradeError
	err := json.Unmarshal(body, &e)
	if err != nil {
		return QuestradeError{
			Code:       -999,
			Message:    "Error unmarshalling error message from Questrade",
			StatusCode: res.StatusCode,
			Endpoint:   res.Request.URL.String(),
		}
	}

	e.StatusCode = res.StatusCode
	e.Endpoint = res.Request.URL.String()

	return e
}

func (q QuestradeError) Error() string {
	return fmt.Sprintf("HTTP %d - %s [%d] - %s", q.StatusCode, q.Endpoint, q.Code, q.Message)
}
