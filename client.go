// Package qapi is a light wrapper for the Questrade REST API, written in Go.
//
// Please note this is not an official API wrapper, and is not endorsed by Questrade. Please see
// http://www.questrade.com/api/home for official documentation.
package qapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// A client is the structure that will be used to consume the API
// endpoints. It holds the login credentials, http client/transport,
// rate limit information, and the login session timer.
type Client struct {
	Credentials        LoginCredentials
	SessionTimer       *time.Timer
	RateLimitRemaining int
	RateLimitReset     time.Time
	httpClient         *http.Client
	transport          *http.Transport
}

// Send an HTTP GET request, and return the processed response
func (c *Client) get(endpoint string, out interface{}, query url.Values) error {
	req, err := http.NewRequest("GET", c.Credentials.ApiServer+endpoint+query.Encode(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.Credentials.authHeader())

	res, err := c.httpClient.Do(req)

	err = c.processResponse(res, out)
	if err != nil {
		return err
	}
	return nil
}

// Format the message body, send an HTTP POST request, and return the processed response
func (c *Client) post(endpoint string, out interface{}, body interface{}) error {
	// Attempt to marshall the body as JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Credentials.ApiServer+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.Credentials.authHeader())

	res, err := c.httpClient.Do(req)

	err = c.processResponse(res, out)
	if err != nil {
		return err
	}
	return nil
}

// processResponse takes the body of an HTTP response, and either returns
// the error code, or unmarshalls the JSON response, extracts
// rate limit info, and places it into the object
// output parameter. This function closes the response body after reading it.
func (c *Client) processResponse(res *http.Response, out interface{}) error {
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return newQuestradeError(res, body)
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}

	reset, _ := strconv.Atoi(res.Header.Get("X-RateLimit-Reset"))
	c.RateLimitReset = time.Unix(int64(reset), 0)
	c.RateLimitRemaining, _ = strconv.Atoi(res.Header.Get("X-RateLimit-Remaining"))

	return nil
}

// Login takes the refresh token from the client login credentials
// and exchanges it for an access token. Returns a timer that
// expires when the login session is over. If the practice flag is
// true, then the client will log into the practice server.
func (c *Client) Login(practice bool) error {
	login := loginServerURL
	if practice {
		login = practiceLoginServerURL
	}

	vars := url.Values{"grant_type": {"refresh_token"}, "refresh_token": {c.Credentials.RefreshToken}}
	res, err := c.httpClient.PostForm(login+"token", vars)

	if err != nil {
		return err
	}

	err = c.processResponse(res, &c.Credentials)
	if err != nil {
		return err
	}

	c.SessionTimer = time.NewTimer(time.Duration(c.Credentials.ExpiresIn) * time.Second)

	return nil
}

// RevokeAuth revokes authorization of the refresh token
// NOTE - You will have to create another manual authorization
// on the Questrade website to use an application again.
func (c *Client) RevokeAuth() error {
	vars := url.Values{"token": {c.Credentials.AccessToken}}

	res, err := c.httpClient.PostForm(loginServerURL+"revoke", vars)

	// Even though the user may still be logged in if there was an error
	// I'm going to set the login info to nil anyways
	c.Credentials = LoginCredentials{}

	defer res.Body.Close()

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Logout Error [%s]: %s", res.Status, body)
	}

	return nil
}

// GetServerTime retrieves the current time on Questrade's server
func (c *Client) GetServerTime() (time.Time, error) {
	t := struct {
		Time time.Time `json:"time"`
	}{}

	err := c.get("v1/time", &t, url.Values{})
	if err != nil {
		return time.Time{}, err
	}

	return t.Time, nil
}

// GetAccounts returns the logged-in User ID, and a list of accounts
// belonging to that user.
func (c *Client) GetAccounts() (int, []Account, error) {
	list := struct {
		UserID   int       `json:"userId"`
		Accounts []Account `json:"accounts"`
	}{}

	err := c.get("v1/accounts", &list, url.Values{})
	if err != nil {
		return 0, []Account{}, err
	}

	return list.UserID, list.Accounts, nil
}

// GetBalances returns the balances for the account with the specified account number
func (c *Client) GetBalances(number string) (AccountBalances, error) {
	bal := AccountBalances{}

	err := c.get("v1/accounts/"+number+"/balances", &bal, url.Values{})
	if err != nil {
		return AccountBalances{}, err
	}

	return bal, nil
}

// GetExecutions returns the number of executions for a given account between the start and end times
// If the times are zero-value, then the API will default the start and end times to the beginning
// and end of the current day.
func (c *Client) GetExecutions(number string, start time.Time, end time.Time) ([]Execution, error) {
	// Format the times if they are not zero-values
	params := url.Values{}
	if !start.Equal(time.Time{}) {
		params.Add("startTime", start.Format(time.RFC3339))
	}

	if !end.Equal(time.Time{}) {
		params.Add("endTime", end.Format(time.RFC3339))
	}

	exec := struct {
		Executions []Execution `json:"executions"`
	}{}

	err := c.get("v1/accounts/"+number+"/executions?", &exec, params)

	if err != nil {
		return []Execution{}, err
	}

	return exec.Executions, nil
}

// GetOrders returns orders for a specified account. Will return results based on the start and
// end times, and the order state. Use GetOrdersByID() to retrieve individual order details.
// If the times are zero-value, then the API will default the start and end times to the beginning
// and end of the current day.
// TODO - Verify order state enumeration in accordance with API docs
// See: http://www.questrade.com/api/documentation/rest-operations/account-calls/accounts-id-orders
func (c *Client) GetOrders(number string, start time.Time, end time.Time, state string) ([]Order, error) {
	// Format the times if they are not zero-values
	params := url.Values{}
	if !start.Equal(time.Time{}) {
		params.Add("startTime", start.Format(time.RFC3339))
	}

	if !end.Equal(time.Time{}) {
		params.Add("endTime", end.Format(time.RFC3339))
	}

	params.Add("stateFilter", state)

	o := struct {
		Orders []Order `json:"orders"`
	}{}

	err := c.get("v1/accounts/"+number+"/orders?", &o, params)
	if err != nil {
		return []Order{}, err
	}

	return o.Orders, nil
}

// GetOrdersByID returns the orders specified by the list of OrderID's
func (c *Client) GetOrdersByID(number string, orderIds ...int) ([]Order, error) {
	idStr := ""
	for k, v := range orderIds {
		idStr += strconv.Itoa(v)
		if k < len(orderIds)-1 {
			idStr += ","
		}
	}

	params := url.Values{}
	params.Add("ids", idStr)

	o := struct {
		Orders []Order `json:"orders"`
	}{}

	err := c.get("v1/accounts/"+number+"/orders?", &o, params)
	if err != nil {
		return []Order{}, err
	}

	return o.Orders, nil
}

// GetSymbols returns detailed symbol information for the given symbol ID's
func (c *Client) GetSymbols(ids ...int) ([]Symbol, error) {
	idStr := ""
	for k, v := range ids {
		idStr += strconv.Itoa(v)
		if k < len(ids)-1 {
			idStr += ","
		}
	}

	params := url.Values{}
	params.Add("ids", idStr)

	s := struct {
		Symbols []Symbol `json:"symbols"`
	}{}

	err := c.get("v1/symbols?", &s, params)
	if err != nil {
		return []Symbol{}, err
	}

	return s.Symbols, nil
}

// SearchSymbols returns symbol search matches for a symbol prefix, at a given offset from the
// beginning of the search results.
func (c *Client) SearchSymbols(prefix string, offset int) ([]SymbolSearchResult, error) {
	params := url.Values{}
	params.Add("prefix", prefix)
	params.Add("offset", strconv.Itoa(offset))

	s := struct {
		Symbols []SymbolSearchResult `json:"symbols"`
	}{}

	err := c.get("v1/symbols/search?", &s, params)
	if err != nil {
		return []SymbolSearchResult{}, err
	}

	return s.Symbols, nil
}

// GetOptionChain Retrieves an option chain for a particular underlying symbol.
// TODO - More comprehensive tests - perhaps I should learn what an option chain is?
func (c *Client) GetOptionChain(id int) ([]OptionChain, error) {
	o := struct {
		Options []OptionChain `json:"options"`
	}{}

	err := c.get("v1/symbols/"+strconv.Itoa(id)+"/options", &o, url.Values{})
	if err != nil {
		return []OptionChain{}, err
	}
	return o.Options, nil
}

// GetMarkets retrieves information about supported markets
func (c *Client) GetMarkets() ([]Market, error) {
	m := struct {
		Markets []Market `json:"markets"`
	}{}

	err := c.get("v1/markets", &m, url.Values{})
	if err != nil {
		return []Market{}, err
	}
	return m.Markets, nil
}

// GetQuote retrieves a single Level 1 market data quote for a single symbol
// TODO - Test
func (c *Client) GetQuote(id int) (Quote, error) {
	idStr := strconv.Itoa(id)

	q := struct {
		Quotes []Quote `json:"quotes"`
	}{}
	//var q2 json.RawMessage

	err := c.get("v1/markets/quotes/"+idStr, &q, url.Values{})
	if err != nil {
		return Quote{}, err
	}

	//fmt.Println("RESULT:", string(q2))

	if len(q.Quotes) != 1 {
		return Quote{}, errors.New("error: Could not retrieve quotes")
	}
	return q.Quotes[0], nil
}

func buildIdString(ids []int) string {
	idStr := ""
	for k, v := range ids {
		idStr += strconv.Itoa(v)
		if k < len(ids)-1 {
			idStr += ","
		}
	}

	return idStr
}

// GetQuotes retrieves a single Level 1 market data quote for many symbols
// TODO - Test
func (c *Client) GetQuotes(ids []int) ([]Quote, error) {
	params := url.Values{}
	params.Add("ids", buildIdString(ids))

	q := struct {
		Quotes []Quote `json:"quotes"`
	}{}

	err := c.get("v1/markets/quotes?", &q, params)
	if err != nil {
		return []Quote{}, err
	}

	return q.Quotes, nil
}

// GetOrderImpact calculates the impact that a given order will have on an
// account without placing it.
// See: http://www.questrade.com/api/documentation/rest-operations/order-calls/accounts-id-orders-impact
func (c *Client) GetOrderImpact(req OrderRequest) (OrderImpact, error) {
	// Construct the endpoint - will be different if the impact is being calculated on
	// an order that already exists
	endpoint := fmt.Sprintf("v1/accounts/%s/orders/", req.AccountID)
	if req.OrderID != 0 {
		endpoint += fmt.Sprintf("%d/", req.OrderID)
	}
	endpoint += "impact"

	var impact OrderImpact
	err := c.post(endpoint, &impact, req)
	if err != nil {
		return OrderImpact{}, err
	}

	return impact, nil
}

// PlaceOrder submits an order request, or an update to an existing order to Questrade
// See: http://www.questrade.com/api/documentation/rest-operations/order-calls/accounts-id-orders
func (c *Client) PlaceOrder(req OrderRequest) ([]Order, error) {
	// Construct the endpoint
	endpoint := fmt.Sprintf("v1/accounts/%s/orders/", req.AccountID)
	if req.OrderID != 0 {
		endpoint += fmt.Sprintf("%d", req.OrderID)
	}

	res := struct {
		OrderID int     `json:"orderId"`
		Orders  []Order `json:"orders"`
	}{}

	err := c.post(endpoint, &res, req)
	if err != nil {
		return []Order{}, err
	}
	return res.Orders, nil
}

// DeleteOrder - Sends a delete request for the specified order
// See: http://www.questrade.com/api/documentation/rest-operations/order-calls/accounts-id-orders-orderid
func (c *Client) DeleteOrder(acctNum string, orderID int) error {
	endpoint := fmt.Sprintf("accounts/%s/orders/%d", acctNum, orderID)
	req, err := http.NewRequest("DELETE", c.Credentials.ApiServer+endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.Credentials.authHeader())

	res, err := c.httpClient.Do(req)

	out := struct {
		OrderID int `json:"OrderID"`
	}{}

	err = c.processResponse(res, out)
	if err != nil {
		return err
	}
	return nil
}

// GetCandles retrieves historical market data between the start and end dates,
// in the given data granularity.
// See: http://www.questrade.com/api/documentation/rest-operations/market-calls/markets-candles-id
func (c *Client) GetCandles(id int, start time.Time, end time.Time, interval string) ([]Candlestick, error) {
	params := url.Values{}
	params.Add("startTime", start.Format(time.RFC3339))
	params.Add("endTime", end.Format(time.RFC3339))
	params.Add("interval", interval)

	r := struct {
		Candles []Candlestick `json:"candles"`
	}{}

	err := c.get("v1/markets/candles/"+strconv.Itoa(id)+"?", &r, params)
	if err != nil {
		return []Candlestick{}, err
	}
	return r.Candles, nil
}

// Gets the port number to which notifications will be streamed
// See: http://www.questrade.com/api/documentation/streaming
func (c *Client) GetNotificationStreamPort(useWebSocket bool) (string, error) {
	mode := "RawSocket"
	if useWebSocket {
		mode = "WebSocket"
	}

	p := struct {
		Port int `json:"streamPort"`
	}{}

	err := c.get("v1/notifications?mode="+mode, &p, url.Values{})
	if err != nil {
		return "", err
	}

	return strconv.Itoa(p.Port), nil
}

// Gets the port number to which quotes will be streamed
// See: http://www.questrade.com/api/documentation/streaming
func (c *Client) GetQuoteStreamPort(useWebSocket bool, ids []int) (string, error) {
	mode := "RawSocket"
	if useWebSocket {
		mode = "WebSocket"
	}

	params := url.Values{}
	params.Add("mode", mode)
	params.Add("stream", "true")
	params.Add("ids", buildIdString(ids))

	p := struct {
		Port int `json:"streamPort"`
	}{}

	err := c.get("v1/markets/quotes?", &p, params)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(p.Port), nil
}

func (c *Client) GetWebSocketConnection(port string) (*websocket.Conn, error) {
	apiServer := c.Credentials.ApiServer[8 : len(c.Credentials.ApiServer)-1]
	conn, _, err := websocket.DefaultDialer.Dial("wss://"+apiServer+":"+port, nil)
	if err != nil {
		return nil, errors.Wrap(err, "WebSocket failed to connect:\n")
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(c.Credentials.AccessToken))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send access token:\n")
	}

	s := struct {
		Success bool `json:"success"`
	}{}

	err = conn.ReadJSON(&s)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read server response:\n")
	}

	return conn, nil
}

// NewClient is the factory function for clients - takes a refresh token and logs into
// either the practice or live server.
func NewClient(refreshToken string, practice bool) (*Client, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new client
	c := &Client{
		Credentials: LoginCredentials{
			RefreshToken: refreshToken,
		},
		httpClient: client,
	}

	err := c.Login(practice)
	if err != nil {
		return nil, err
	}

	return c, nil
}
