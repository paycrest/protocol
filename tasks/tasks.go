package tasks

import (
	"context"

	"github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
	"github.com/paycrest/paycrest-protocol/services"
	"github.com/paycrest/paycrest-protocol/utils/logger"
)

// ContinueIndexing continues indexing
func ContinueIndexing() error {
	ctx := context.Background()
	indexerService := services.NewIndexerService(nil)

	// Start ERC20 transfer indexing
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
		receiveAddress := receiveAddress

		go indexerService.RunIndexERC20Transfer(ctx, receiveAddress)
	}

	// Start indexing on-chain payment order deposits
	// TODO: query networks based on the development environment: prod == mainnet, sandbox == testnet
	networks, err := database.GetClient().Network.Query().All(ctx)
	if err != nil {
		return err
	}

	for _, network := range networks {
		go func(network *ent.Network) {
			err := indexerService.IndexOrderDeposits(ctx, nil, network)
			if err != nil {
				logger.Errorf("process order deposits task => %v\n", err)
			}
		}(network)
	}

	return nil
}

// ProcessOrders processes orders to the on-chain escrow
func ProcessOrders() error {
	ctx := context.Background()

	orders, err := database.GetClient().PaymentOrder.
		Query().
		Where(
			paymentorder.StatusEQ(paymentorder.StatusInitiated),
			paymentorder.HasReceiveAddressWith(
				receiveaddress.StatusEQ(receiveaddress.StatusUsed),
			),
		).
		All(ctx)
	if err != nil {
		return err
	}

	go func() {
		for _, order := range orders {
			orderService := services.NewOrderService()
			order := order

			err := orderService.CreateOrder(ctx, nil, order.ID)
			if err != nil {
				logger.Errorf("process orders task => %v\n", err)
			}
		}
	}()

	return nil
}
