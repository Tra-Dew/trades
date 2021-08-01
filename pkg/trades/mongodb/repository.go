package mongodb

import (
	"context"
	"time"

	"github.com/d-leme/tradew-trades/pkg/trades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repositoryMongoDB struct {
	collection *mongo.Collection
}

// NewRepository ...
func NewRepository(client *mongo.Client, database string) trades.Repository {
	repository := &repositoryMongoDB{client.Database(database).Collection("trades")}
	repository.createIndex()

	return repository
}

// Insert ...
func (repository *repositoryMongoDB) Insert(ctx context.Context, trade *trades.TradeOffer) error {

	_, err := repository.collection.InsertOne(ctx, trade)

	return err
}

// Update ...
func (repository *repositoryMongoDB) Update(ctx context.Context, trade *trades.TradeOffer) error {

	filter := bson.M{"_id": trade.ID}

	_, err := repository.collection.UpdateOne(ctx, filter, trade)

	return err
}

// Get ...
func (repository *repositoryMongoDB) Get(ctx context.Context, userID string, req *trades.GetTradesOffers) (*trades.ResultTradeOffers, error) {

	if req.PageSize < 1 {
		req.PageSize = 10
	}

	result := new(trades.ResultTradeOffers)
	result.Trades = []*trades.TradeOffer{}

	filter := bson.M{}
	if req.Token != nil {
		filter["_id"] = bson.M{"$gt": req.Token}
	}

	cursor, err := repository.collection.Find(
		ctx,
		filter,
		options.Find().SetSort(bson.M{"_id": 1}).SetLimit(req.PageSize),
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &result.Trades)
	if err != nil {
		return nil, err
	}

	if lastItem := result.Trades[len(result.Trades)-1]; lastItem != nil {
		result.Token = lastItem.ID
	}

	return result, nil
}

// GetByID ...
func (repository *repositoryMongoDB) GetByID(ctx context.Context, userID, id string) (*trades.TradeOffer, error) {
	var result *trades.TradeOffer

	err := repository.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repository *repositoryMongoDB) createIndex() {
	_, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

}
