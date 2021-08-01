package trades_test

import (
	"context"
	"errors"
	"testing"

	"github.com/d-leme/tradew-trades/pkg/core"
	"github.com/d-leme/tradew-trades/pkg/trades"
	"github.com/d-leme/tradew-trades/pkg/trades/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type serviceTestSuite struct {
	suite.Suite
	assert           *assert.Assertions
	ctx              context.Context
	repository       *mock.RepositoryMock
	service          trades.Service
	inventoryService *mock.InventoryServiceMock
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (s *serviceTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()
}

func (s *serviceTestSuite) SetupTest() {
	s.repository = mock.NewRepository().(*mock.RepositoryMock)
	s.inventoryService = mock.NewInventoryService().(*mock.InventoryServiceMock)
	s.service = trades.NewService(s.repository, s.inventoryService)
}

func (s *serviceTestSuite) TestCreate() {

	userID := uuid.NewString()
	wantedItemsOwnerID := uuid.NewString()
	correlationID := uuid.NewString()

	req := &trades.CreateTradeOfferRequest{
		WantedItemsOwnerID: wantedItemsOwnerID,
		OfferedItems: []*trades.ItemModel{
			{
				ID:       uuid.NewString(),
				Quantity: 5,
			},
		},
		WantedItems: []*trades.ItemModel{
			{
				ID:       uuid.NewString(),
				Quantity: 5,
			},
		},
	}

	s.repository.On("Insert").Return(nil)
	s.repository.On("Update").Return(nil)
	s.inventoryService.On("LockItems").Return(nil)

	res, err := s.service.Create(s.ctx, userID, correlationID, req)

	s.assert.NoError(err)
	s.assert.NotNil(res)
	s.assert.NotEmpty(res.ID)

	s.repository.AssertNumberOfCalls(s.T(), "Insert", 1)
	s.repository.AssertNumberOfCalls(s.T(), "Update", 1)
	s.inventoryService.AssertNumberOfCalls(s.T(), "LockItems", 1)
}

func (s *serviceTestSuite) TestCreateLockFailed() {
	userID := uuid.NewString()
	wantedItemsOwnerID := uuid.NewString()
	correlationID := uuid.NewString()

	req := &trades.CreateTradeOfferRequest{
		WantedItemsOwnerID: wantedItemsOwnerID,
		OfferedItems: []*trades.ItemModel{
			{
				ID:       uuid.NewString(),
				Quantity: 5,
			},
		},
		WantedItems: []*trades.ItemModel{
			{
				ID:       uuid.NewString(),
				Quantity: 5,
			},
		},
	}

	s.repository.On("Insert").Return(nil)
	s.repository.On("Update").Return(nil)
	s.inventoryService.On("LockItems").Return(errors.New("invalid-wanted-items"))

	res, err := s.service.Create(s.ctx, userID, correlationID, req)

	s.assert.ErrorIs(core.ErrLockFailed, err)
	s.assert.Nil(res)

	s.repository.AssertNumberOfCalls(s.T(), "Insert", 1)
	s.repository.AssertNumberOfCalls(s.T(), "Update", 1)
	s.inventoryService.AssertNumberOfCalls(s.T(), "LockItems", 1)
}

func (s *serviceTestSuite) TestAccept() {
	correlationID := uuid.NewString()

	trade := &trades.TradeOffer{
		ID:                 uuid.NewString(),
		OwnerID:            uuid.NewString(),
		WantedItemsOwnerID: uuid.NewString(),
		Status:             trades.TradePending,
		OfferedItems: []*trades.Item{
			{
				ID:       uuid.NewString(),
				Quantity: 1,
			},
		},
		WantedItems: []*trades.Item{
			{
				ID:       uuid.NewString(),
				Quantity: 2,
			},
		},
	}

	s.repository.On("GetByID", trade.ID).Return(trade)
	s.repository.On("Update").Return(nil)
	s.inventoryService.On("TradesItems").Return(nil)

	err := s.service.Accept(s.ctx, trade.OwnerID, correlationID, trade.ID)

	s.assert.NoError(err)

	s.repository.AssertNumberOfCalls(s.T(), "GetByID", 1)
	s.repository.AssertNumberOfCalls(s.T(), "Update", 2)
	s.inventoryService.AssertNumberOfCalls(s.T(), "TradesItems", 1)
}

func (s *serviceTestSuite) TestAcceptInvalidStatus() {
	correlationID := uuid.NewString()

	trade := &trades.TradeOffer{
		ID:                 uuid.NewString(),
		OwnerID:            uuid.NewString(),
		WantedItemsOwnerID: uuid.NewString(),
		Status:             trades.TradeCreated,
		OfferedItems: []*trades.Item{
			{
				ID:       uuid.NewString(),
				Quantity: 1,
			},
		},
		WantedItems: []*trades.Item{
			{
				ID:       uuid.NewString(),
				Quantity: 2,
			},
		},
	}

	s.repository.On("GetByID", trade.ID).Return(trade)

	err := s.service.Accept(s.ctx, trade.OwnerID, correlationID, trade.ID)

	s.assert.ErrorIs(core.ErrTradeInvalidStatus, err)

	s.repository.AssertNumberOfCalls(s.T(), "GetByID", 1)
}

func (s *serviceTestSuite) TestAcceptTradeItemsFailed() {
	correlationID := uuid.NewString()

	trade := &trades.TradeOffer{
		ID:                 uuid.NewString(),
		OwnerID:            uuid.NewString(),
		WantedItemsOwnerID: uuid.NewString(),
		Status:             trades.TradePending,
		OfferedItems: []*trades.Item{
			{
				ID:       uuid.NewString(),
				Quantity: 1,
			},
		},
		WantedItems: []*trades.Item{
			{
				ID:       uuid.NewString(),
				Quantity: 2,
			},
		},
	}

	s.repository.On("GetByID", trade.ID).Return(trade)
	s.repository.On("Update").Return(nil)
	s.inventoryService.On("TradesItems").Return(errors.New("unable-to-trade-items"))

	err := s.service.Accept(s.ctx, trade.OwnerID, correlationID, trade.ID)

	s.assert.ErrorIs(core.ErrItemsTradeFailed, err)

	s.repository.AssertNumberOfCalls(s.T(), "GetByID", 1)
	s.repository.AssertNumberOfCalls(s.T(), "Update", 2)
	s.inventoryService.AssertNumberOfCalls(s.T(), "TradesItems", 1)
}
