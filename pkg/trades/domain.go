package trades

import (
	"context"
	"time"

	"github.com/Tra-Dew/trades/pkg/core"
)

// TradeStatus ...
type TradeStatus string

const (
	// TradeDispatchLockOfferedItems ...
	TradeDispatchLockOfferedItems TradeStatus = "DispatchLockOfferedItemsPending"

	// TradeLockOfferedItemsPending ...
	TradeLockOfferedItemsPending TradeStatus = "LockOfferedItemsPending"

	// TradeCreated ...
	TradeCreated TradeStatus = "Created"

	// TradeDispatchTradeAcceptedPending ...
	TradeDispatchTradeAcceptedPending TradeStatus = "DispatchTradeAcceptedPending"

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
	Create(ctx context.Context, trade *TradeOffer) error
	Update(ctx context.Context, userID string, trade *TradeOffer) error
	Get(ctx context.Context, userID string, req *GetTradesOffers) (*ResultTradeOffers, error)
	GetByID(ctx context.Context, userID, id string) (*Item, error)
}

// Service ...
type Service interface {
	Create(ctx context.Context, userID, correlationID string, req *CreateTradeOfferRequest) (*CreateTradeOfferResponse, error)
	Accept(ctx context.Context, userID, correlationID, id string) error
	Cancel(ctx context.Context, userID, correlationID, id string) error
	Refuse(ctx context.Context, userID, correlationID, id string) error
	Get(ctx context.Context, userID, correlationID string, req *GetTradeOffersRequest) (*GetTradeOffersResponse, error)
	GetByID(ctx context.Context, userID, id string) (*GetTradeOfferResponse, error)
}

// NewTradeOffer ...
func NewTradeOffer(id, ownerID string, status string, offeredItems, wantedItems []*Item) (*TradeOffer, error) {

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
		Status:       TradeStatus(status),
		OfferedItems: offeredItems,
		WantedItems:  wantedItems,
		CreatedAt:    time.Now(),
	}, nil
}
