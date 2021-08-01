package mock

import (
	"context"

	"github.com/d-leme/tradew-trades/pkg/trades/external/inventory"
	"github.com/stretchr/testify/mock"
)

// InventoryServiceMock ...
type InventoryServiceMock struct {
	mock.Mock
}

// NewInventoryService ...
func NewInventoryService() inventory.Service {
	return &InventoryServiceMock{}
}

// LockItems ...
func (r *InventoryServiceMock) LockItems(ctx context.Context, req *inventory.LockItemsRequest) error {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(error)
	}

	return nil
}

// TradesItems ...
func (r *InventoryServiceMock) TradesItems(ctx context.Context, req *inventory.TradeItemsRequest) error {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(error)
	}

	return nil
}
