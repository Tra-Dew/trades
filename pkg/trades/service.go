package trades

import (
	"context"

	"github.com/d-leme/tradew-trades/pkg/core"
	"github.com/d-leme/tradew-trades/pkg/trades/external/inventory"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

	fields := logrus.Fields{
		"user_id":               userID,
		"wanted_items_owner_id": req.WantedItemsOwnerID,
		"correlation_id":        correlationID,
	}

	offeredItems, err := ToDomain(req.OfferedItems)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error parsing offered items")
		return nil, err
	}

	wantedItems, err := ToDomain(req.WantedItems)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error parsing wanted items")
		return nil, err
	}

	trade, err := NewTradeOffer(uuid.NewString(), userID, req.WantedItemsOwnerID, offeredItems, wantedItems)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error creating new offer")

		return nil, err
	}

	if err := s.repository.Insert(ctx, trade); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error inserting offer")
		return nil, err
	}

	fields["trade_id"] = trade.ID
	logrus.WithFields(fields).Info("new trade created")

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
		logrus.WithError(err).WithFields(fields).Error("error locking items")

		trade.UpdateStatus(TradeError)

		if err := s.repository.Update(ctx, trade); err != nil {
			logrus.WithError(err).WithFields(fields).Error("error updating trade")

			return nil, err
		}

		logrus.WithFields(fields).Info("trade status set to error")

		return nil, core.ErrLockFailed
	}

	trade.UpdateStatus(TradePending)

	if err := s.repository.Update(ctx, trade); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error updating trade")
		return nil, err
	}

	logrus.WithFields(fields).Info("trade status set to pending")

	return &CreateTradeOfferResponse{ID: trade.ID}, nil
}

func (s *service) Accept(ctx context.Context, userID, correlationID, id string) error {

	fields := logrus.Fields{
		"trade_id":       id,
		"user_id":        userID,
		"correlation_id": correlationID,
	}

	trade, err := s.repository.GetByID(ctx, userID, id)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error getting trade")
		return err
	}

	if trade.Status != TradePending {
		logrus.
			WithError(core.ErrTradeInvalidStatus).
			WithFields(fields).
			Error("tried to accept trade that was in an invalid state")

		return core.ErrTradeInvalidStatus
	}

	trade.UpdateStatus(TradeAccepted)

	if err := s.repository.Update(ctx, trade); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error updating trade")
		return err
	}

	logrus.WithFields(fields).Info("trade status set to accepted")

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

		trade.UpdateStatus(TradeError)

		if err := s.repository.Update(ctx, trade); err != nil {
			logrus.WithError(err).WithFields(fields).Error("error updating trade")
			return err
		}

		logrus.WithFields(fields).Info("trade status set to error")

		return core.ErrItemsTradeFailed
	}

	trade.UpdateStatus(TradeCompleted)

	if err := s.repository.Update(ctx, trade); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error updating trade")
		return err
	}

	logrus.WithFields(fields).Info("trade status set to completed")

	return nil
}

func (s *service) Get(ctx context.Context, userID string, req *GetTradeOffersRequest) (*GetTradeOffersResponse, error) {

	res, err := s.repository.Get(ctx, userID, &GetTradesOffers{
		Token:    req.Token,
		PageSize: req.PageSize,
	})

	if err != nil {
		logrus.
			WithError(err).
			WithFields(logrus.Fields{
				"user_id": userID,
				"token":   req.Token,
			}).
			Error("error getting trades")
		return nil, err
	}

	return ParseGetTradeOffersResponse(res), nil
}

func (s *service) GetByID(ctx context.Context, userID, id string) (*GetTradeOfferResponse, error) {

	trade, err := s.repository.GetByID(ctx, userID, id)
	if err != nil {
		logrus.
			WithError(err).
			WithFields(logrus.Fields{
				"trade_id": id,
				"user_id":  userID,
			}).
			Error("error getting trades")
		return nil, err
	}

	return &GetTradeOfferResponse{Trade: ParseTradeOffer(trade)}, nil
}
