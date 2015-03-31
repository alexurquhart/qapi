package qapi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type QuestradeError struct {
	Code               int `json:"code",string`
	StatusCode         int
	Message            string `json:"message"`
	Endpoint           string
	OrderId            int     `json:"orderId"`
	Orders             []Order `json:"orders"`
	RateLimitRemaining int
	RateLimitReset     time.Time
}

// TODO - Parse rate limits out of response
func newQuestradeError(res *http.Response, body []byte) QuestradeError {
	reset, _ := strconv.Atoi(res.Header.Get("X-RateLimit-Reset"))
	resetDate := time.Unix(int64(reset), 0)
	remaining, _ := strconv.Atoi(res.Header.Get("X-RateLimit-Remaining"))

	e := QuestradeError{
		StatusCode:         res.StatusCode,
		Endpoint:           res.Request.URL.String(),
		RateLimitRemaining: remaining,
		RateLimitReset:     resetDate,
	}
	return e
}

func (q QuestradeError) Error() string {
	return fmt.Sprintf("%d %s [%d] - %s", q.StatusCode, q.Endpoint, q.Code, q.Message)
}
