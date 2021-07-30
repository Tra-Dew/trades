package inventory

import "context"

// ItemToLock ...
type ItemToLock struct {
	ID       string
	Quantity int64
}

// LockItemsRequest ...
type LockItemsRequest struct {
	LockedBy string
	OwnerID  string
	Items    []*ItemToLock
}

// ItemToTrade ...
type ItemToTrade struct {
	ID       string
	Quantity int64
}

// TradeItemsRequest ...
type TradeItemsRequest struct {
	TradeID            string
	OwnerID            string
	WantedItemsOwnerID string
	OfferedItems       []*ItemToTrade
	WantedItems        []*ItemToTrade
}

// Service ...
type Service interface {
	LockItems(ctx context.Context, req *LockItemsRequest) error
	TradesItems(ctx context.Context, req *TradeItemsRequest) error
}
