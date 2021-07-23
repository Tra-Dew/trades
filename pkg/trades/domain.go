package trades

import (
	"context"
	"time"

	"github.com/Tra-Dew/trades/pkg/core"
)

// TradeStatus ...
type TradeStatus string

const (
	// TradeCreated ...
	TradeCreated TradeStatus = "Created"

	// TradeLockOfferedItemsPending ...
	TradeLockOfferedItemsPending TradeStatus = "LockOfferedItemsPending"

	// TradePending ...
	TradePending TradeStatus = "Pending"

	// TradeAccepted ...
	TradeAccepted TradeStatus = "Accepted"

	// TradeAwaitingItemsTrade ...
	TradeAwaitingItemsTrade TradeStatus = "AwaitingItemsTrade"

	// TradeCompleted ...
	TradeCompleted TradeStatus = "Completed"
)

// Item ...
type Item struct {
	ID       string `bson:"id"`
	Quantity int64  `bson:"quantity"`
}

// TradeOffer ...
type TradeOffer struct {
	ID           string      `bson:"_id"`
	OwnerID      string      `bson:"owner_id"`
	Status       TradeStatus `bson:"status"`
	OfferedItems []*Item     `bson:"offered_items"`
	WantedItems  []*Item     `bson:"wanted_items"`
	CreatedAt    time.Time   `bson:"created_at"`
	UpdatedAt    *time.Time  `bson:"updated_at"`
}

// GetTradesOffers ...
type GetTradesOffers struct {
	Token    *string
	PageSize int64
}

// ResultTradeOffers ...
type ResultTradeOffers struct {
	Trades []*TradeOffer
	Token  string
}

// Repository ...
type Repository interface {
	Insert(ctx context.Context, trade *TradeOffer) error
	Update(ctx context.Context, trade *TradeOffer) error
	UpdateBulk(ctx context.Context, trade []*TradeOffer) error
	Get(ctx context.Context, userID string, req *GetTradesOffers) (*ResultTradeOffers, error)
	GetByID(ctx context.Context, userID string, id string) (*TradeOffer, error)
	GetByIDs(ctx context.Context, ids []string) ([]*TradeOffer, error)
	GetByStatus(ctx context.Context, status TradeStatus) ([]*TradeOffer, error)
}

// Service ...
type Service interface {
	Create(ctx context.Context, userID, correlationID string, req *CreateTradeOfferRequest) (*CreateTradeOfferResponse, error)
	Accept(ctx context.Context, userID, correlationID, id string) error
	Get(ctx context.Context, userID string, req *GetTradeOffersRequest) (*GetTradeOffersResponse, error)
	GetByID(ctx context.Context, userID, id string) (*GetTradeOfferResponse, error)
}

// NewItem ...
func NewItem(id string, quantity int64) (*Item, error) {

	if id == "" {
		return nil, core.ErrValidationFailed
	}

	if quantity < 1 {
		return nil, core.ErrValidationFailed
	}

	return &Item{
		ID:       id,
		Quantity: quantity,
	}, nil
}

// NewTradeOffer ...
func NewTradeOffer(id, ownerID string, offeredItems, wantedItems []*Item) (*TradeOffer, error) {

	if id == "" {
		return nil, core.ErrValidationFailed
	}

	if ownerID == "" {
		return nil, core.ErrValidationFailed
	}

	if len(offeredItems) < 1 {
		return nil, core.ErrValidationFailed
	}

	if len(wantedItems) < 1 {
		return nil, core.ErrValidationFailed
	}

	return &TradeOffer{
		ID:           id,
		OwnerID:      ownerID,
		Status:       TradeCreated,
		OfferedItems: offeredItems,
		WantedItems:  wantedItems,
		CreatedAt:    time.Now(),
	}, nil
}

// UpdateStatus ...
func (trade *TradeOffer) UpdateStatus(status TradeStatus) {
	trade.Status = status

	now := time.Now()
	trade.UpdatedAt = &now
}
