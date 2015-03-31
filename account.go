package qapi

import (
	"time"
)

// Account represents an account associated with the user on behalf
// of which the API client is authorized.
//
// Ref: http://www.questrade.com/api/documentation/rest-operations/account-calls/accounts
type Account struct {
	// Type of the account (e.g., "Cash", "Margin").
	Type string `json":"type""`

	// Eight-digit account number (e.g., "26598145")
	// Stored as a string, it's used for making account-related API calls
	Number string `json:"number"`

	// Status of the account (e.g., Active).
	Status string `json:"status"`

	// Whether this is a primary account for the holder.
	IsPrimary bool `json:"isPrimary"`

	// Whether this account is one that gets billed for various expenses such as inactivity fees, market data, etc.
	IsBilling bool `json:"isBilling"`

	// Type of client holding the account (e.g., "Individual").
	ClientAccountType string `json:"clientAccountType"`
}

// Position belonging to an account
//
// Ref: http://www.questrade.com/api/documentation/rest-operations/account-calls/accounts-id-positions
type Position struct {
	// Position symbol.
	Symbol string `json:"symbol"`

	// Internal symbol identifier
	SymbolID int `json:"symbolId"`

	// Position quantity remaining open.
	OpenQuantity float32 `json:"openQuantity"`

	// Portion of the position that was closed today.
	ClosedQuantity float32 `json:"closedQuantity"`

	// Market value of the position (quantity x price).
	CurrentMarketValue float32 `json:"currentMarketValue"`

	// Current price of the position symbol.
	CurrentPrice float32 `json:"currentPrice"`

	// Average price paid for all executions constituting the position.
	AverageEntryPrice float32 `json:"averageEntryPrice"`

	// Realized profit/loss on this position.
	ClosedPnL float32 `json:"closedPnL"`

	// Unrealized profit/loss on this position.
	OpenPnL float32 `json:"openPnL"`

	// Total cost of the position.
	TotalCost float32 `json:"totalCost"`

	// Designates whether real-time quote was used to compute PnL.
	IsRealTime bool `json:"isRealTime"`

	// Designates whether a symbol is currently undergoing a reorg.
	IsUnderReorg bool `json:"isUnderReorg"`
}

// Balance belonging to an Account
//
// Ref: http://www.questrade.com/api/documentation/rest-operations/account-calls/accounts-id-balances
type Balance struct {

	// Currency of the balance figure(e.g., "USD" or "CAD").
	Currency string `json:"currency"`

	// Balance amount.
	Cash float32 `json:"cash"`

	// Market value of all securities in the account in a given currency.
	MarketValue float32 `json:"marketValue"`

	// Equity as a difference between cash and marketValue properties.
	TotalEquity float32 `json:"totalEquity"`

	// Buying power for that particular currency side of the account.
	BuyingPower float32 `json:"buyingPower"`

	// Maintenance excess for that particular side of the account.
	MaintenanceExcess float32 `json:"maintenanceExcess"`

	// Whether real-time data was used to calculate the above values.
	IsRealTime bool `json:"isRealTime"`
}

// AccountBalances represents per-currency and combined balances for a specified account.
//
// Ref: http://www.questrade.com/api/documentation/rest-operations/account-calls/accounts-id-balances
type AccountBalances struct {
	PerCurrencyBalances    []Balance `json:"perCurrencyBalances"`
	CombinedBalances       []Balance `json:"combinedBalances"`
	SODPerCurrencyBalances []Balance `json:"sodPerCurrencyBalances"`
	SODCombinedBalances    []Balance `json:"sodCombinedBalances"`
}

// Execution belonging to an Account
//
// Ref: http://www.questrade.com/api/documentation/rest-operations/account-calls/accounts-id-executions
type Execution struct {
	// Execution symbol.
	Symbol string `json:"symbol"`

	// Internal symbol identifier
	SymbolID int `json:"symbolId"`

	// Execution quantity.
	Quantity int `json:"quantity"`

	// Client side of the order to which execution belongs.
	Side string `json:"side"`

	// Execution price.
	Price float32 `json:"price"`

	// Internal identifier of the execution.
	ID int `json:"id"`

	// Internal identifier of the order to which the execution belongs.
	OrderID int `json:"orderId"`

	// Internal identifier of the order chain to which the execution belongs.
	OrderChainID int `json:"orderChainId"`

	// Identifier of the execution at the market where it originated.
	ExchangeExecID string `json:"exchangeExecId"`

	// Execution timestamp.
	Timestamp time.Time `json:"timestamp"`

	// Manual notes that may have been entered by Trade Desk staff
	Notes string `json:"notes"`

	// Trading venue where execution originated.
	Venue string `json:"venue"`

	// Execution cost (price x quantity).
	TotalCost float32 `json:"totalCost"`

	// Questrade commission for orders placed with Trade Desk.
	OrderPlacementCommission float32 `json:"orderPlacementCommission"`

	// Questrade commission.
	Commission float32 `json:"commission"`

	// Liquidity fee charged by execution venue.
	ExecutionFee float32 `json:"executionFee"`

	// SEC fee charged on all sales of US securities.
	SecFee float32 `json:"secFee"`

	// Additional execution fee charged by TSX (if applicable).
	CanadianExecutionFee int `json:"canadianExecutionFee"`

	// Internal identifierof the parent order.
	ParentID int `json:"parentId"`
}
