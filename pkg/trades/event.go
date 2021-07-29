package trades

import "time"

// ItemEventModel ...
type ItemEventModel struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

// TradeOfferCreatedEvent ...
type TradeOfferCreatedEvent struct {
	ID                 string            `json:"id"`
	OwnerID            string            `json:"owner_id"`
	WantedItemsOwnerID string            `json:"wanted_items_owner_id"`
	OfferedItems       []*ItemEventModel `json:"offered_items"`
	WantedItems        []*ItemEventModel `json:"wanted_items"`
	CreatedAt          time.Time         `json:"created_at"`
}

// ItemsLockCompletedEvent ...
type ItemsLockCompletedEvent struct {
	LockedBy string `json:"locked_by"`
}

// TradeOfferAcceptedEvent ...
type TradeOfferAcceptedEvent struct {
	ID                 string            `json:"id"`
	OwnerID            string            `json:"owner_id"`
	WantedItemsOwnerID string            `json:"wanted_items_owner_id"`
	OfferedItems       []*ItemEventModel `json:"offered_items"`
	WantedItems        []*ItemEventModel `json:"wanted_items"`
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
		ID:                 trade.ID,
		OwnerID:            trade.OwnerID,
		WantedItemsOwnerID: trade.WantedItemsOwnerID,
		OfferedItems:       ParseItemsToItemSliceModel(trade.OfferedItems),
		WantedItems:        ParseItemsToItemSliceModel(trade.WantedItems),
		CreatedAt:          trade.CreatedAt,
	}
}

// ParseTradeToTradeOfferAcceptedEvent ...
func ParseTradeToTradeOfferAcceptedEvent(trade *TradeOffer) *TradeOfferAcceptedEvent {
	return &TradeOfferAcceptedEvent{
		ID:                 trade.ID,
		OwnerID:            trade.OwnerID,
		WantedItemsOwnerID: trade.WantedItemsOwnerID,
		WantedItems:        ParseItemsToItemSliceModel(trade.WantedItems),
	}
}
