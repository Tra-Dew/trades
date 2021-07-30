package proto

import (
	"context"

	"github.com/d-leme/tradew-trades/pkg/trades/external/inventory"
)

type service struct {
	client InventoryServiceClient
}

// NewService ...
func NewService(client InventoryServiceClient) inventory.Service {
	return &service{
		client: client,
	}
}

func (s *service) LockItems(ctx context.Context, req *inventory.LockItemsRequest) error {

	items := make([]*ItemToLock, len(req.Items))

	for i, item := range req.Items {
		items[i] = &ItemToLock{
			Id:       item.ID,
			Quantity: item.Quantity,
		}
	}

	protoReq := &LockItemsRequest{
		LockedBy: req.LockedBy,
		OwnerID:  req.OwnerID,
		Items:    items,
	}

	if _, err := s.client.LockItems(context.Background(), protoReq); err != nil {
		return err
	}

	return nil
}

func (s *service) TradesItems(ctx context.Context, req *inventory.TradeItemsRequest) error {

	offeredItems := make([]*ItemToTrade, len(req.OfferedItems))
	for i, item := range req.OfferedItems {
		offeredItems[i] = &ItemToTrade{
			Id:       item.ID,
			Quantity: item.Quantity,
		}
	}

	wantedItems := make([]*ItemToTrade, len(req.WantedItems))
	for i, item := range req.WantedItems {
		wantedItems[i] = &ItemToTrade{
			Id:       item.ID,
			Quantity: item.Quantity,
		}
	}

	protoReq := &TradeItemsRequest{
		TradeID:            req.TradeID,
		OwnerID:            req.OwnerID,
		WantedItemsOwnerID: req.WantedItemsOwnerID,
		OfferedItems:       offeredItems,
		WantedItems:        wantedItems,
	}

	if _, err := s.client.TradeItems(context.Background(), protoReq); err != nil {
		return err
	}

	return nil
}
