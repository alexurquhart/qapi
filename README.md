#qapi

qapi is an abstraction over the [Questrade REST API](http://www.questrade.com/api/documentation/getting-started), written in Go.

##Logging In
You will need to sign up for a practice account, or create an API key for your trading account. You can do so at
http://www.questrade.com/api/home

To login, you will need to provide the refresh token generated on the Questrade API admin site
During login, the refresh token is exchanged for an access token, which qapi uses to make authenticated
API calls. Upon successful login, Questrade returns a replacement refresh token, which is necessary to login
again after the current session expires. See http://www.questrade.com/api/documentation/security for more info.

##Usage
```go
// Create a new client
client, err := qapi.NewClient(“< REFRESH TOKEN >”)

// Get accounts
userId, accts, err := client.GetAccounts()

// Get balances for the first account
balances, err := client.GetBalances(accts[0].Number)

// Print the market value of the combined balances
for b := range balances.CombinedBalances {
    fmt.Printf(“Market Value: $%.2f %s\n”, b.MarketValue, b.Currency)
}

// To get a quote the API uses internal symbol ID’s.
results, err := client.SearchSymbols(“AAPL”, 0)

var symId int
for r := range results {
    if r.Symbol == “AAPL” && r.ListingExchange == “NASDAQ” {
        symId = r.SymbolID
        break
    }
}

// Get a real-time quote - qapi supports getting quotes of multiple symbols with GetQuotes()
quote, err := client.GetQuote(symId)
fmt.Printf(“AAPL Bid Price: $%.2f\n”, quote.BidPrice)

// We’re done with the application forever - deauthorize the API key
client.RevokeAuth()
```

##Disclaimer
**NOTE** - This library is not endorsed or supported by Questrade in any way, shape or form. This library is released under the MIT License.
