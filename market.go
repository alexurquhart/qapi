package qapi

import (
	"time"
)

// Detailed information about a symbol
// Ref: http://www.questrade.com/api/documentation/rest-operations/market-calls/symbols-id
type Symbol struct {

	// Symbol that follows Questrade symbology (e.g., "TD.TO").
	Symbol string `json:"symbol"`

	// Symbol identifier
	SymbolID int `json:"symbolId"`

	// Closing trade price from the previous trading day.
	PrevDayClosePrice float32 `json:"prevDayClosePrice"`

	// 52-week high price.
	HighPrice52 float32 `json:"highPrice52"`

	// 52-week low price.
	LowPrice52 float32 `json:"lowPrice52"`

	// Average trading volume over trailing 3 months.
	AverageVol3Months int `json:"averageVol3Months"`

	// Average trading volume over trailing 20 days.
	AverageVol20Days int `json:"averageVol20Days"`

	// Total number of shares outstanding.
	OutstandingShares int `json:"outstandingShares"`

	// Trailing 12-month earnings per share.
	EPS float32 `json:"eps"`

	// Trailing 12-month price to earnings ratio.
	PE float32 `json:"pe"`

	// Dividend amount per share.
	Dividend float32 `json:"dividend"`

	// Dividend yield (dividend / prevDayClosePrice).
	Yield float32 `json:"yield"`

	// Dividend ex-date.
	ExDate *time.Time `json:"exDate"`

	// Market capitalization (outstandingShares * prevDayClosePrice).
	MarketCap float32 `json:"marketCap"`

	// Option type (e.g., "Call").
	OptionType string `json:"optionType"`

	// Option duration type (e.g., "Weekly").
	OptionDurationType string `json:"optionDurationType"`

	// Option root symbol (e.g., "MSFT").
	OptionRoot string `json:"optionRoot"`

	// Option contract deliverables.
	OptionContractDeliverables OptionContractDeliverables `json:"optionContractDeliverables"`

	// Option exercise style (e.g., "American").
	OptionExerciseType string `json:"optionExerciseType"`

	// Primary listing exchange.
	ListingExchange string `json:"listingExchange"`

	// Symbol description (e.g., "Microsoft Corp.").
	Description string `json:"description"`

	// Security type (e.g., "Stock").
	SecurityType string `json:"securityType"`

	// Option expiry date.
	OptionExpiryDate *time.Time `json:"optionExpiryDate"`

	// Dividend declaration date.
	DividendDate *time.Time `json:"dividendDate"`

	// Option strike price.
	OptionStrikePrice float32 `json:"optionStrikePrice"`

	// Indicates whether the symbol is actively listed.
	IsQuotable bool `json:"isQuotable"`

	// Indicates whether the symbol is an underlying option.
	HasOptions bool `json:"hasOptions"`

	// String Currency code (follows ISO format).
	Currency string `json:"currency"`

	// List of MinTickData records.
	MinTicks []MinTickData `json:"minTicks"`
}

type UnderlyingMultiplierPair struct {
	Multiplier         int    `json:"multiplier"`
	UnderlyingSymbol   string `json:"underlyingSymbol"`
	UnderlyingSymbolID string `json:"underlyingSymbolId"`
}

type OptionContractDeliverables struct {
	Underlyings []UnderlyingMultiplierPair `json:"underlyings"`
	CashInLieu  float32                    `json:"cashInLieu"`
}

type MinTickData struct {
	Pivot   float32 `json:"pivot"`
	MinTick float32 `json:"minTick"`
}

// Symbol information retreived from search results
// Ref: http://www.questrade.com/api/documentation/rest-operations/market-calls/symbols-search
type SymbolSearchResult struct {
	Symbol          string `json:"symbol"`
	SymbolID        int    `json:"symbolId"`
	Description     string `json:"description"`
	SecurityType    string `json:"securityType"`
	ListingExchange string `json:"listingExchange"`
	IsQuotable      bool   `json:"isQuotable"`
	IsTradable      bool   `json:"isTradable"`
	Currency        string `json:"currency"`
}

type ChainPerStrikePrice struct {
	// Option strike price.
	StrikePrice float32 `json:"strikePrice"`

	// Internal identifier of the call option symbol.
	CallSymbolID int `json:"callSymbolId"`

	// Internal identifier of the put option symbol.
	PutSymbolID int `json:"putSymbolId"`
}

type ChainPerRoot struct {
	// Option root symbol.
	Root string `json:"root"`

	// Slice of ChainPerStrikePrice elements.
	ChainPerStrikePrice []ChainPerStrikePrice `json:"chainPerStrikePrice"`

	// NOTE - This value appears in example API response, but is not documented
	Multiplier int `json:"multiplier"`
}

// Option Chain
// Ref: www.questrade.com/api/documentation/rest-operations/market-calls/symbols-id-options
type OptionChain struct {
	// Option expiry date.
	ExpiryDate time.Time `json:"expiryDate"`

	// Description of the underlying option.
	Description string `json:"description"`

	// Primary listing exchange.
	ListingExchange string `json:"listingExchange"`

	// Option exercise style (e.g., "American").
	OptionExerciseType string `json:"optionExerciseType"`

	// Slice of ChainPerRoot elements
	ChainPerRoot []ChainPerRoot `json:"chainPerRoot"`
}

// Market represents information about supported markets
type Market struct {
	// Market name.
	Name string `json:"name"`

	// List of trading venue codes.
	TradingVenues []string `json:"tradingVenues"`

	// Default trading venue code.
	DefaultTradingVenue string `json:"defaultTradingVenue"`

	// List of primary order route codes.
	PrimaryOrderRoutes []string `json:"primaryOrderRoutes"`

	// List of secondary order route codes.
	SecondaryOrderRoutes []string `json:"secondaryOrderRoutes"`

	// List of level 1 market data feed codes.
	Level1Feeds []string `json:"level1Feeds"`

	// List of level 2 market data feed codes.
	Level2Feeds []string `json:"level2Feeds"`

	// Pre-market opening time for current trading date.
	ExtendedStartTime time.Time `json:"extendedStartTime"`

	// Regular market opening time for current trading date.
	StartTime time.Time `json:"startTime"`

	// Regular market closing time for current trading date.
	EndTime time.Time `json:"endTime"`

	// Extended market closing time for current trading date.
	ExtendedEndTime time.Time `json:"extendedEndTime"`

	// Currency code (ISO format).
	Currency string `json:"currency"`

	// Number of snap quotes that the user can retrieve from a market.
	SnapQuotesLimit int `json:"snapQuotesLimit"`
}

// Quote represents a Lvl 1 market data quote for a particular symbol
type Quote struct {
	// Symbol name following Questrade’s symbology.
	Symbol int `json:"symbol"`

	// Internal symbol identifier.
	SymbolID string `json:"symbolId"`

	// Market tier.
	Tier string `json:"tier"`

	// Bid price.
	BidPrice float32 `json:"bidPrice"`

	// Bid quantity.
	BidSize int `json:"bidSize"`

	// Ask price.
	AskPrice float32 `json:"askPrice"`

	// Ask quantity.
	AskSize int `json:"askSize"`

	// Price of the last trade during regular trade hours.
	LastTradeTrHrs float32 `json:"lastTradeTrHrs"`

	// Price of the last trade.
	LastTradePrice float32 `json:"lastTradePrice"`

	// Quantity of the last trade.
	LastTradeSize int `json:"lastTradeSize"`

	// Trade direction.
	LastTradeTick string `json:"lastTradeTick"`

	// Volume.
	Volume int `json:"volume"`

	// Opening trade price.
	OpenPrice float32 `json:"openPrice"`

	// Daily high price.
	HighPrice float32 `json:"highPrice"`

	// Daily low price.
	LowPrice float32 `json:"lowPrice"`

	// Whether a quote is delayed (true) or real-time.
	Delay bool `json:"delay"`

	// Whether trading in the symbol is currently halted.
	IsHalted bool `json:"isHalted"`
}

// Candlestick represents historical market data in the form of OHLC candlesticks
// for a specified symbol.
type Candlestick struct {
	// Candlestick start timestamp
	Start time.Time `json:"start"`

	// Candlestick end timestamp
	End time.Time `json:"end"`

	// Opening price.
	Open float32 `json:"open"`

	// High price.
	High float32 `json:"high"`

	// Low price.
	Low float32 `json:"low"`

	// Closing price.
	Close float32 `json:"close"`

	// Trading volume.
	Volume int `json:"volume"`
}
