package cmd

import (
	"context"

	"github.com/Tra-Dew/trades/pkg/core"
	"github.com/Tra-Dew/trades/pkg/trades"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Container ...
type Container struct {
	Settings *core.Settings

	Authenticate *core.Authenticate

	MongoClient *mongo.Client

	SNS *session.Session
	SQS *session.Session

	TradeRepository trades.Repository
	TradeService    trades.Service
	TradeController trades.Controller
}

// NewContainer creates new instace of Container
func NewContainer(settings *core.Settings) *Container {

	container := new(Container)

	container.Settings = settings

	container.SQS = session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(settings.SQS.Region),
		Endpoint: aws.String(settings.SQS.Endpoint),
	}))

	container.SNS = session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(settings.SNS.Region),
		Endpoint: aws.String(settings.SNS.Endpoint),
	}))

	container.MongoClient = connectMongoDB(settings.MongoDB)

	container.Authenticate = core.NewAuthenticate(settings.JWT.Secret)

	// container.InventoryRepository = mongodb.NewRepository(container.MongoClient, settings.MongoDB.Database)
	// container.InventoryController = inventory.NewController(container.Authenticate, container.InventoryRepository)

	return container
}

// Controllers maps all routes and exposes them
func (c *Container) Controllers() []core.Controller {
	return []core.Controller{
		&c.TradeController,
	}
}

// Close terminates every opened resource
func (c *Container) Close() {
	c.MongoClient.Disconnect(context.Background())
}

func connectMongoDB(conf *core.MongoDBConfig) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(conf.ConnectionString))

	if err != nil {
		logrus.
			WithError(err).
			Fatal("error connecting to MongoDB")
	}

	client.Connect(context.Background())

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		logrus.
			WithError(err).
			Fatal("error pinging MongoDB")
	}

	return client
}
