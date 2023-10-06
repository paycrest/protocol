package tasks

import (
	"context"
	"time"

	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/lockpaymentorder"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
	"github.com/paycrest/paycrest-protocol/services"
	"github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/types"
	"github.com/paycrest/paycrest-protocol/utils/logger"
)

// ContinueIndexing continues indexing
func ContinueIndexing() error {
	ctx := context.Background()
	indexerService := services.NewIndexerService(nil)

	// Start ERC20 transfer indexing
	addresses, err := storage.GetClient().ReceiveAddress.
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

	// Start indexing on-chain payment order deposits and settlements
	// TODO: query networks based on the development environment: prod == mainnet, sandbox == testnet
	networks, err := storage.GetClient().Network.Query().All(ctx)
	if err != nil {
		return err
	}

	for _, network := range networks {
		go func(network *ent.Network) {
			err := indexerService.IndexOrderDeposits(ctx, nil, network)
			if err != nil {
				logger.Errorf("process order deposits task => %v\n", err)
			}

			err = indexerService.IndexOrderSettlements(ctx, nil, network)
			if err != nil {
				logger.Errorf("process order settlements task => %v\n", err)
			}
		}(network)
	}

	return nil
}

// ProcessOrders processes orders to the on-chain escrow
func ProcessOrders() error {
	ctx := context.Background()

	orders, err := storage.GetClient().PaymentOrder.
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

// SubscribeToRedisKeyspaceEvents subscribes to redis keyspace events according to redis.conf settings
func SubscribeToRedisKeyspaceEvents() {
	ctx := context.Background()

	// Handle expired or deleted order request key events
	orderRequest := storage.RedisClient.PSubscribe(
		ctx,
		"__keyevent@0__:expired:order_request_*",
		"__keyevent@0__:del:order_request_*",
	)
	orderRequestChan := orderRequest.Channel()

	go services.NewPriorityQueueService().ReassignStaleOrderRequest(ctx, orderRequestChan)
}

// ProcessUnfulfilledLockOrders reassigns lockOrder not unfulfilled within a time frame.
func ProcessUnfulfilledLockOrders() error {
	ctx := context.Background()

	lockOrders, err := storage.GetClient().LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusProcessing),
			lockpaymentorder.UpdatedAtLTE(time.Now().Add(-services.OrderConf.ReceiveAddressValidity*time.Minute)),
		).
		All(ctx)
	if err != nil {
		return err
	}

	go func() {
		for _, order := range lockOrders {

			lockPaymentOrder := types.LockPaymentOrderFields{
				Token:             order.Edges.Token,
				OrderID:           order.OrderID,
				Amount:            order.Amount,
				Rate:              order.Rate,
				BlockNumber:       order.BlockNumber,
				Institution:       order.Institution,
				AccountIdentifier: order.AccountIdentifier,
				AccountName:       order.AccountName,
				ProviderID:        order.Edges.Provider.ID,
				ProvisionBucket:   order.Edges.ProvisionBucket,
			}

			err := services.NewPriorityQueueService().AssignLockPaymentOrder(ctx, lockPaymentOrder)
			if err != nil {
				logger.Errorf("task reassign unfulfilled lock order with id: %s => %v\n", order.OrderID, err)
			}
		}
	}()

	return nil
}
