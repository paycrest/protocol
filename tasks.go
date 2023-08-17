package main

import (
	"context"

	"github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
	"github.com/paycrest/paycrest-protocol/services"
	"github.com/paycrest/paycrest-protocol/utils/logger"
)

// ContinueIndexing continues indexing
func ContinueIndexing() error {
	ctx := context.Background()

	addresses, err := database.GetClient().ReceiveAddress.
		Query().
		Where(
			receiveaddress.Or(
				receiveaddress.StatusEQ(receiveaddress.StatusUnused),
				receiveaddress.StatusEQ(receiveaddress.StatusPartial),
			),
		).All(ctx)
	if err != nil {
		return err
	}

	for _, receiveAddress := range addresses {
		indexerService := services.NewIndexerService(nil)
		receiveAddress := receiveAddress

		go indexerService.RunIndexERC20Transfer(ctx, receiveAddress)
	}

	return nil
}

// ProcessOrders processes orders to the on-chain escrow
func ProcessOrders() error {
	ctx := context.Background()

	orders, err := database.GetClient().PaymentOrder.
		Query().
		Where(
			paymentorder.And(
				paymentorder.StatusEQ(paymentorder.StatusInitiated),
				paymentorder.HasReceiveAddressWith(
					receiveaddress.StatusEQ(receiveaddress.StatusUsed),
				),
			),
		).
		All(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		orderService := services.NewOrderService()
		order := order

		go func() {
			err := orderService.CreateOrder(ctx, nil, order.ID)
			if err != nil {
				logger.Errorf("error: %v", err)
			}
		}()
	}

	return nil
}
