package trades

import "time"

// ItemModel ...
type ItemModel struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

// TradeOfferModel ...
type TradeOfferModel struct {
	ID           string       `json:"_id"`
	Status       string       `json:"status"`
	OfferedItems []*ItemModel `json:"offered_items"`
	WantedItems  []*ItemModel `json:"wanted_items"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    *time.Time   `json:"updated_at"`
}

// CreateTradeOfferRequest ...
type CreateTradeOfferRequest struct {
	OfferedItems []*ItemModel `json:"offered_items"`
	WantedItems  []*ItemModel `json:"wanted_items"`
}

// CreateTradeOfferResponse ...
type CreateTradeOfferResponse struct {
	ID string `json:"id"`
}

// GetTradeOffersRequest ...
type GetTradeOffersRequest struct {
	Token    *string `form:"token"`
	PageSize int64   `form:"page_size"`
}

// GetTradeOffersResponse ...
type GetTradeOffersResponse struct {
	Trades []*TradeOfferModel `json:"trades"`
	Token  string             `json:"token"`
}

// GetTradeOfferResponse ...
type GetTradeOfferResponse struct {
	Trade *TradeOfferModel `json:"trade"`
}

// ParseItem ...
func ParseItem(item *Item) *ItemModel {
	return &ItemModel{
		ID:       item.ID,
		Quantity: item.Quantity,
	}
}

// ParseItemSlice ...
func ParseItemSlice(s []*Item) []*ItemModel {
	items := make([]*ItemModel, len(s))

	for i, item := range s {
		items[i] = ParseItem(item)
	}
	return items
}

// ParseTradeOffer ...
func ParseTradeOffer(trade *TradeOffer) *TradeOfferModel {
	return &TradeOfferModel{
		ID:           trade.ID,
		Status:       string(trade.Status),
		OfferedItems: ParseItemSlice(trade.OfferedItems),
		WantedItems:  ParseItemSlice(trade.WantedItems),
		CreatedAt:    trade.CreatedAt,
		UpdatedAt:    trade.UpdatedAt,
	}
}

// ParseTradeOfferSlice ...
func ParseTradeOfferSlice(s []*TradeOffer) []*TradeOfferModel {
	trades := make([]*TradeOfferModel, len(s))

	for i, trade := range s {
		trades[i] = ParseTradeOffer(trade)
	}

	return trades
}

// ParseGetTradeOffersResponse ...
func ParseGetTradeOffersResponse(res *ResultTradeOffers) *GetTradeOffersResponse {
	return &GetTradeOffersResponse{
		Token:  res.Token,
		Trades: ParseTradeOfferSlice(res.Trades),
	}
}

// ToDomain ...
func ToDomain(models []*ItemModel) ([]*Item, error) {
	items := make([]*Item, len(models))

	for i, item := range models {

		newItem, err := NewItem(item.ID, item.Quantity)
		if err != nil {
			return nil, err
		}

		items[i] = newItem
	}

	return items, nil
}
