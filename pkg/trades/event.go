package trades

import "time"

// ItemEventModel ...
type ItemEventModel struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

// TradeOfferCreatedEvent ...
type TradeOfferCreatedEvent struct {
	ID           string            `json:"id"`
	OwnerID      string            `json:"owner_id"`
	Status       string            `json:"status"`
	OfferedItems []*ItemEventModel `json:"offered_items"`
	WantedItems  []*ItemEventModel `json:"wanted_items"`
	CreatedAt    time.Time         `json:"created_at"`
}

// ItemLockCompletedEvent ...
type ItemLockCompletedEvent struct {
	ID        string    `json:"id"`
	LockedBy  string    `json:"locked_by"`
	Quantity  int64     `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ItemsLockCompletedEvent ...
type ItemsLockCompletedEvent struct {
	Items []*ItemLockCompletedEvent `json:"items"`
}

// TradeOfferAcceptedEvent ...
type TradeOfferAcceptedEvent struct {
	ID          string            `json:"id"`
	OwnerID     string            `json:"owner_id"`
	WantedItems []*ItemEventModel `json:"wanted_items"`
}

// ItemsTradeCompletedEvent ...
type ItemsTradeCompletedEvent struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
}

// ParseItemsToItemSliceModel ...
func ParseItemsToItemSliceModel(items []*Item) []*ItemEventModel {
	events := make([]*ItemEventModel, len(items))

	for i, item := range items {
		events[i] = &ItemEventModel{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	return events
}

// ParseTradeToTradeOfferCreatedEvent ...
func ParseTradeToTradeOfferCreatedEvent(trade *TradeOffer) *TradeOfferCreatedEvent {
	return &TradeOfferCreatedEvent{
		ID:           trade.ID,
		OwnerID:      trade.OwnerID,
		Status:       string(trade.Status),
		OfferedItems: ParseItemsToItemSliceModel(trade.OfferedItems),
		WantedItems:  ParseItemsToItemSliceModel(trade.WantedItems),
		CreatedAt:    trade.CreatedAt,
	}
}

// ParseTradeToTradeOfferAcceptedEvent ...
func ParseTradeToTradeOfferAcceptedEvent(trade *TradeOffer) *TradeOfferAcceptedEvent {
	return &TradeOfferAcceptedEvent{
		ID:          trade.ID,
		OwnerID:     trade.OwnerID,
		WantedItems: ParseItemsToItemSliceModel(trade.WantedItems),
	}
}
