package qapi

import (
	"time"
)

// TODO - Populate this struct
type OrderLeg struct {
}

// Ref: http://www.questrade.com/api/documentation/rest-operations/account-calls/accounts-id-orders
type Order struct {
	// Internal order identifier.
	ID int `json:"id"`

	// Symbol that follows Questrade symbology (e.g., "TD.TO").
	Symbol string `json:"symbol"`

	// Internal symbol identifier.
	SymbolID int `json:"symbolId"`

	// Total quantity of the order.
	TotalQuantity int `json:"totalQuantity"`

	// Unfilled portion of the order quantity.
	OpenQuantity int `json:"openQuantity"`

	// Filled portion of the order quantity.
	FilledQuantity int `json:"filledQuantity"`

	// Unfilled portion of the order quantity after cancellation.
	CanceledQuantity int `json:"canceledQuantity"`

	// Client view of the order side (e.g., "Buy-To-Open").
	Side string `json:"side"`

	// Order price type (e.g., "Market").
	OrderType string `json:"orderType"`

	// Limit price.
	LimitPrice float32 `json:"limitPrice"`

	// Stop price.
	StopPrice float32 `json:"stopPrice"`

	// Specifies all-or-none special instruction.
	IsAllOrNone bool `json:"isAllOrNone"`

	// Specifies Anonymous special instruction.
	IsAnonymous bool `json:"isAnonymous"`

	// Specifies Iceberg special instruction.
	IcebergQuantity int `json:"icebergQuantity"`

	// Specifies Minimum special instruction.
	MinQuantity int `json:"minQuantity"`

	// Average price of all executions received for this order.
	AvgExecPrice float32 `json:"avgExecPrice"`

	// Price of the last execution received for the order in question.
	LastExecPrice float32 `json:"lastExecPrice"`

	// See enumerations for all allowed values
	Source string `json:"source"`

	// See Order Time In Force section for all allowed values.
	TimeInForce string `json:"timeInForce"`

	// Good-Till-Date marker and date parameter
	GtdDate *time.Time `json:"gtdDate",string`

	// See Order State section for all allowed values.
	State string `json:"state"`

	// Human readable order rejection reason message.
	ClientReasonStr string `json:"clientReasonStr"`

	// Internal identifier of a chain to which the order belongs.
	ChainID int `json:"chainId"`

	// Order creation time.
	CreationTime *time.Time `json:"creationTime",string`

	// Time of the last update.
	UpdateTime *time.Time `json:"updateTime",string`

	// Notes that may have been manually added by Questrade staff.
	Notes string `json:"notes"`

	// See enumerations for all allowed values
	PrimaryRoute string `json:"primaryRoute"`

	// See enumerations for all allowed values
	SecondaryRoute string `json:"secondaryRoute"`

	// Order route name.
	OrderRoute string `json:"orderRoute"`

	// Venue where non-marketable portion of the order was booked.
	VenueHoldingOrder string `json:"venueHoldingOrder"`

	// Total commission amount charged for this order.
	CommissionCharged float32 `json:"commissionCharged"`

	// Identifier assigned to this order by exchange where it was routed.
	ExchangeOrderID string `json:"exchangeOrderId"`

	// Whether user that placed the order is a significant shareholder.
	IsSignificantShareholder bool `json:"isSignificantShareholder"`

	// Whether user that placed the order is an insider.
	IsInsider bool `json:"isInsider"`

	// Whether limit offset is specified in dollars (vs. percent).
	IsLimitOffsetInDollar bool `json:"isLimitOffsetInDollar"`

	// Internal identifier of user that placed the order.
	UserID int `json:"userId"`

	// Commission for placing the order via the Trade Desk over the phone.
	PlacementCommission float32 `json:"placementCommission"`

	// List of OrderLeg elements.
	Legs []OrderLeg `json:"legs"`

	// Multi-leg strategy to which the order belongs.
	StrategyType string `json:"strategyType"`

	// Stop price at which order was triggered.
	TriggerStopPrice float32 `json:"triggerStopPrice"`

	// Internal identifier of the order group.
	OrderGroupID int `json:"orderGroupId"`

	// Bracket Order class. Primary, Profit or Loss.
	OrderClass string `json:"orderClass"`
}
