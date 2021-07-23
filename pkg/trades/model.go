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
