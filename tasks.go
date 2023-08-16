package main

import (
	"context"

	"github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
	"github.com/paycrest/paycrest-protocol/services"
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
