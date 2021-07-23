package cmd

import (
	"context"

	"github.com/Tra-Dew/trades/pkg/core"
	"github.com/Tra-Dew/trades/pkg/trades"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DispatchTradeOfferCreated ...
func DispatchTradeOfferCreated(command *cobra.Command, args []string) {
	settings := new(core.Settings)

	err := core.FromYAML(command.Flag("settings").Value.String(), settings)
	if err != nil {
		logrus.
			WithError(err).
			Fatal("unable to parse settings, shutting down...")
	}

	ctx := context.Background()
	container := NewContainer(settings)

	pendingTrades, err := container.TradeRepository.GetByStatus(ctx, trades.TradeCreated)
	if err != nil {
		logrus.
			WithError(err).
			Error("error while getting trades by status")
		return
	}

	lenItems := len(pendingTrades)

	logrus.Infof("%d new trades to publish", lenItems)

	if lenItems < 1 {
		return
	}

	for _, trade := range pendingTrades {
		fields := logrus.Fields{"trade_id": trade.ID, "owner_id": trade.OwnerID}

		event := trades.ParseTradeToTradeOfferCreatedEvent(trade)

		messageID, err := container.Producer.Publish(settings.Events.TradeCreated, event)

		if err != nil {
			logrus.
				WithError(err).
				WithFields(fields).
				Error("error while dispatching message")
			continue
		}

		fields["message_id"] = messageID

		logrus.
			WithFields(fields).
			Info("dipached event")

		trade.UpdateStatus(trades.TradeLockOfferedItemsPending)

		if err := container.TradeRepository.Update(ctx, trade); err != nil {
			logrus.
				WithError(err).
				WithFields(fields).
				Error("error while updating trade")
			continue
		}
	}

	logrus.Info("worker complete")
}
