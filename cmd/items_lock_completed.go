package cmd

import (
	"context"
	"reflect"

	"github.com/d-leme/tradew-trades/pkg/core"
	"github.com/d-leme/tradew-trades/pkg/trades"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ItemsLockCompleted ...
func ItemsLockCompleted(command *cobra.Command, args []string) {
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
		core.WithSubscriberID(settings.Events.ItemsLockCompleted),
		core.WithTopicID(settings.Events.ItemsLockCompleted),
		core.WithMaxRetries(3),
		core.WithType(reflect.TypeOf(trades.ItemsLockCompletedEvent{})),
		core.WithHandler(func(payload interface{}) error {
			message := payload.(*trades.ItemsLockCompletedEvent)

			logrus.Info("processing received event")

			ctx := context.Background()

			itemMap := map[string]*trades.ItemLockCompletedEvent{}
			ids := []string{}
			for _, item := range message.Items {
				_, exists := itemMap[item.ID]

				if !exists {
					itemMap[item.ID] = item
					ids = append(ids, item.ID)
				}
			}

			fields := logrus.Fields{"ids": ids}

			pendingTrades, err := container.TradeRepository.GetByIDs(ctx, ids)
			if err != nil {
				logrus.
					WithError(err).
					WithFields(fields).
					Error("error while getting trades by ids")
				return err
			}

			for _, trade := range pendingTrades {
				trade.UpdateStatus(trades.TradePending)
			}

			if err := container.TradeRepository.UpdateBulk(ctx, pendingTrades); err != nil {
				logrus.
					WithError(err).
					WithFields(fields).
					Error("error while updating trades")
				return err
			}

			logrus.
				WithFields(fields).
				Info("trades set to pending successfully")

			return nil
		}))

	if err := consumer.Run(); err != nil {
		logrus.
			WithError(err).
			Error("shutting down with error")
	}
}
