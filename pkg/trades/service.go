package trades

import (
	"context"

	"github.com/d-leme/tradew-trades/pkg/core"
	"github.com/d-leme/tradew-trades/pkg/trades/external/inventory"
	"github.com/google/uuid"
)

type service struct {
	repository       Repository
	inventoryService inventory.Service
}

// NewService ...
func NewService(repository Repository, inventoryService inventory.Service) Service {
	return &service{
		repository:       repository,
		inventoryService: inventoryService,
	}
}

func (s *service) Create(
	ctx context.Context,
	userID, correlationID string,
	req *CreateTradeOfferRequest,
) (*CreateTradeOfferResponse, error) {

	offeredItems, err := ToDomain(req.OfferedItems)
	if err != nil {
		return nil, err
	}

	wantedItems, err := ToDomain(req.WantedItems)
	if err != nil {
		return nil, err
	}

	trade, err := NewTradeOffer(uuid.NewString(), userID, req.WantedItemsOwnerID, offeredItems, wantedItems)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Insert(ctx, trade); err != nil {
		return nil, err
	}

	lockItemsReq := &inventory.LockItemsRequest{
		LockedBy:           trade.ID,
		OwnerID:            trade.OwnerID,
		WantedItemsOwnerID: trade.WantedItemsOwnerID,
		OfferedItems:       make([]*inventory.ItemToLock, len(trade.OfferedItems)),
		WantedItems:        make([]*inventory.ItemToLock, len(trade.WantedItems)),
	}

	for i, item := range trade.OfferedItems {
		lockItemsReq.OfferedItems[i] = &inventory.ItemToLock{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	for i, item := range trade.WantedItems {
		lockItemsReq.WantedItems[i] = &inventory.ItemToLock{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	if err := s.inventoryService.LockItems(ctx, lockItemsReq); err != nil {
		return nil, err
	}

	trade.UpdateStatus(TradePending)

	if s.repository.Update(ctx, trade); err != nil {
		return nil, err
	}

	return &CreateTradeOfferResponse{ID: trade.ID}, nil
}

func (s *service) Accept(ctx context.Context, userID, correlationID, id string) error {
	trade, err := s.repository.GetByID(ctx, userID, id)
	if err != nil {
		return err
	}

	if trade.Status != TradeCreated {
		return core.ErrTradeCompleteInvalidStatus
	}

	trade.UpdateStatus(TradeAccepted)

	if err := s.repository.Update(ctx, trade); err != nil {
		return err
	}

	tradeItemsReq := &inventory.TradeItemsRequest{
		TradeID:            trade.ID,
		OwnerID:            trade.OwnerID,
		WantedItemsOwnerID: trade.WantedItemsOwnerID,
		OfferedItems:       make([]*inventory.ItemToTrade, len(trade.OfferedItems)),
		WantedItems:        make([]*inventory.ItemToTrade, len(trade.WantedItems)),
	}

	for i, item := range trade.OfferedItems {
		tradeItemsReq.OfferedItems[i] = &inventory.ItemToTrade{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	for i, item := range trade.WantedItems {
		tradeItemsReq.WantedItems[i] = &inventory.ItemToTrade{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	if err := s.inventoryService.TradesItems(ctx, tradeItemsReq); err != nil {
		return err
	}

	return nil
}

func (s *service) Get(ctx context.Context, userID string, req *GetTradeOffersRequest) (*GetTradeOffersResponse, error) {

	res, err := s.repository.Get(ctx, userID, &GetTradesOffers{
		Token:    req.Token,
		PageSize: req.PageSize,
	})

	if err != nil {
		return nil, err
	}

	return ParseGetTradeOffersResponse(res), nil
}

func (s *service) GetByID(ctx context.Context, userID, id string) (*GetTradeOfferResponse, error) {

	trade, err := s.repository.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	return &GetTradeOfferResponse{Trade: ParseTradeOffer(trade)}, nil
}
