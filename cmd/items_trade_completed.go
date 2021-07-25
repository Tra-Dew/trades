package cmd

import (
	"context"
	"reflect"

	"github.com/d-leme/tradew-trades/pkg/core"
	"github.com/d-leme/tradew-trades/pkg/trades"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ItemsTradeCompleted ...
func ItemsTradeCompleted(command *cobra.Command, args []string) {
	settings := new(core.Settings)

	err := core.FromYAML(command.Flag("settings").Value.String(), settings)
	if err != nil {
		logrus.
			WithError(err).
			Fatal("unable to parse settings, shutting down...")
	}

	container := NewContainer(settings)

	consumer := core.NewMessageBrokerSubscriber(
		core.WithSessionSNS(container.SNS),
		core.WithSessionSQS(container.SQS),
		core.WithSubscriberID(settings.Events.ItemsTradeCompleted),
		core.WithTopicID(settings.Events.ItemsTradeCompleted),
		core.WithMaxRetries(3),
		core.WithType(reflect.TypeOf(trades.ItemsTradeCompletedEvent{})),
		core.WithHandler(func(payload interface{}) error {
			message := payload.(*trades.ItemsTradeCompletedEvent)

			logrus.Info("processing received event")

			ctx := context.Background()

			fields := logrus.Fields{"trade_id": message.ID, "owner_id": message.OwnerID}

			trade, err := container.TradeRepository.GetByID(ctx, message.OwnerID, message.ID)
			if err != nil {
				logrus.
					WithError(err).
					WithFields(fields).
					Error("error while getting trades by ids")
				return err
			}

			trade.UpdateStatus(trades.TradeCompleted)

			if err := container.TradeRepository.Update(ctx, trade); err != nil {
				logrus.
					WithError(err).
					WithFields(fields).
					Error("error while updating trades")
				return err
			}

			logrus.
				WithFields(fields).
				Info("trade completed successfully")

			return nil
		}))

	if err := consumer.Run(); err != nil {
		logrus.
			WithError(err).
			Error("shutting down with error")
	}
}
