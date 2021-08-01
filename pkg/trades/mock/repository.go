package mock

import (
	"context"

	"github.com/d-leme/tradew-trades/pkg/trades"
	"github.com/stretchr/testify/mock"
)

// RepositoryMock ...
type RepositoryMock struct {
	mock.Mock
}

// NewRepository ...
func NewRepository() trades.Repository {
	return &RepositoryMock{}
}

// Insert ...
func (r *RepositoryMock) Insert(ctx context.Context, trade *trades.TradeOffer) error {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(error)
	}

	return nil
}

// Update ...
func (r *RepositoryMock) Update(ctx context.Context, trade *trades.TradeOffer) error {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(error)
	}

	return nil
}

// Get ...
func (r *RepositoryMock) Get(ctx context.Context, userID string, req *trades.GetTradesOffers) (*trades.ResultTradeOffers, error) {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(*trades.ResultTradeOffers), nil
	}

	arg1 := args.Get(1)

	return nil, arg1.(error)
}

// GetByID ...
func (r *RepositoryMock) GetByID(ctx context.Context, userID string, id string) (*trades.TradeOffer, error) {
	args := r.Mock.Called(id)

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(*trades.TradeOffer), nil
	}

	arg1 := args.Get(1)

	return nil, arg1.(error)
}
