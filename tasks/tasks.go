package tasks

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/ent/providerordertoken"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
	"github.com/paycrest/paycrest-protocol/services"
	"github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/shopspring/decimal"
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

// ComputeMarketRate computes the market price for fiat currencies
func ComputeMarketRate() error {
	ctx := context.Background()

	// Fetch all fiat currencies
	currencies, err := storage.Client.FiatCurrency.
		Query().
		All(ctx)
	if err != nil {
		return err
	}

	priorityQueue := services.NewPriorityQueueService()

	for _, currency := range currencies {
		go func(currency *ent.FiatCurrency) {
			// Fetch rates from token configs with fixed conversion rate
			tokenConfigs, err := storage.Client.ProviderOrderToken.
				Query().
				Where(
					providerordertoken.SymbolEQ(providerordertoken.SymbolUSDT),
					providerordertoken.ConversionRateTypeEQ(providerordertoken.ConversionRateTypeFixed),
				).
				All(ctx)
			if err != nil {
				logger.Errorf("compute market price task => %v\n", err)
				return
			}

			var rates []decimal.Decimal
			for _, tokenConfig := range tokenConfigs {
				rates = append(rates, tokenConfig.FixedConversionRate)
			}

			// Calculate median
			median := utils.Median(rates)

			// Update currency with median rate
			currency, err = storage.Client.FiatCurrency.
				UpdateOneID(currency.ID).
				SetMarketRate(median).
				Save(ctx)
			if err != nil {
				logger.Errorf("compute market price task => %v\n", err)
			}

			// Create default bucket for currency
			priorityQueue.CreatePriorityQueueForDefaultBucket(ctx, currency)
		}(currency)
	}

	return nil
}

// StartCronJobs starts cron jobs
func StartCronJobs() {
	scheduler := gocron.NewScheduler(time.UTC)

	// Compute market rate four times a day - starting at 6AM
	_, err := scheduler.Cron("0 6,12,18,0 * * *").Do(ComputeMarketRate)
	if err != nil {
		logger.Errorf("failed to schedule compute market rate task => %v\n", err)
	}

	// Start scheduler
	scheduler.StartAsync()
}
