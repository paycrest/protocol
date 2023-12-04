package tasks

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/providerordertoken"
	"github.com/paycrest/protocol/ent/receiveaddress"
	"github.com/paycrest/protocol/ent/webhookretryattempt"
	"github.com/paycrest/protocol/services"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/logger"
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
		Where(fiatcurrency.IsEnabledEQ(true)).
		All(ctx)
	if err != nil {
		return err
	}

	priorityQueue := services.NewPriorityQueueService()

	// Fetch stable coin rate from third-party API Binance (USDT)
	resp, err := utils.MakeJSONRequest(
		ctx,
		"GET",
		"https://api.binance.com/api/v3/ticker/price?symbol=USDTNGN",
		nil,
		nil,
	)
	if err != nil {
		logger.Errorf("failed to fetch third-party rate => %v\n", err)
		return err
	}

	var externalRate decimal.Decimal
	if resp != nil {
		externalRate, _ = decimal.NewFromString(resp["price"].(string))
	}

	for _, currency := range currencies {
		// Fetch rates from token configs with fixed conversion rate
		tokenConfigs, err := storage.Client.ProviderOrderToken.
			Query().
			Where(
				providerordertoken.SymbolEQ("USDT"),
				providerordertoken.ConversionRateTypeEQ(providerordertoken.ConversionRateTypeFixed),
			).
			All(ctx)
		if err != nil {
			logger.Errorf("compute market price task => %v\n", err)
		}

		var rates []decimal.Decimal
		for _, tokenConfig := range tokenConfigs {
			rates = append(rates, tokenConfig.FixedConversionRate)
		}

		// Calculate median
		median := utils.Median(rates)

		// Check the median rate against the external rate to ensure it's not too far off
		allowedDeviation := decimal.NewFromFloat(0.005) // 0.5%
		if externalRate.Cmp(decimal.Zero) != 0 {
			if median.LessThan(externalRate.Mul(decimal.NewFromFloat(1).Sub(allowedDeviation))) ||
				median.GreaterThan(externalRate.Mul(decimal.NewFromFloat(1).Add(allowedDeviation))) {
				median = externalRate
			}
		}

		// Update currency with median rate
		currency, err = storage.Client.FiatCurrency.
			UpdateOneID(currency.ID).
			SetMarketRate(median).
			Save(ctx)
		if err != nil {
			logger.Errorf("compute market price task => %v\n", err)
			return err
		}

		// Create default bucket for currency
		priorityQueue.CreatePriorityQueueForDefaultBucket(ctx, currency)
	}

	return nil
}

// Retry failed webhook notifications
func RetryFailedWebhookNotifications() error {
	ctx := context.Background()

	// Fetch failed webhook notifications that are due for retry
	attempts, err := storage.Client.WebhookRetryAttempt.
		Query().
		Where(
			webhookretryattempt.StatusEQ(webhookretryattempt.StatusFailed),
			webhookretryattempt.NextRetryTimeLTE(time.Now()),
		).
		All(ctx)
	if err != nil {
		logger.Errorf("RetryFailedWebhookNotifications: %v\n", err)
		return err
	}

	baseDelay := 2 * time.Minute
	maxCumulativeTime := 24 * time.Hour

	for _, attempt := range attempts {
		// Send the webhook notification
		_, err = utils.MakeJSONRequest(
			ctx,
			"POST",
			attempt.WebhookURL,
			attempt.Payload,
			map[string]string{
				"X-Paycrest-Signature": attempt.Signature,
			},
		)
		if err != nil {
			// Webhook notification failed
			// Update attempt with next retry time
			attemptNumber := attempt.AttemptNumber + 1
			delay := baseDelay * time.Duration(math.Pow(2, float64(attemptNumber-1)))

			nextRetryTime := time.Now().Add(delay)

			attemptUpdate := attempt.Update()

			attemptUpdate.
				AddAttemptNumber(1).
				SetNextRetryTime(nextRetryTime)

			// Set status to expired if cumulative time is greater than 24 hours
			if nextRetryTime.Sub(attempt.CreatedAt.Add(-baseDelay)) > maxCumulativeTime {
				attemptUpdate.SetStatus(webhookretryattempt.StatusExpired)
			}

			_, err := attemptUpdate.Save(ctx)
			if err != nil {
				logger.Errorf("RetryFailedWebhookNotifications: %v\n", err)
			}

			continue
		}

		// Webhook notification was successful
		_, err := attempt.Update().
			SetStatus(webhookretryattempt.StatusSuccess).
			Save(ctx)
		if err != nil {
			logger.Errorf("RetryFailedWebhookNotifications: %v\n", err)
		}
	}

	return nil
}

// StartCronJobs starts cron jobs
func StartCronJobs() {
	conf := config.OrderConfig()
	scheduler := gocron.NewScheduler(time.UTC)

	// Compute market rate four times a day - starting at 6AM
	_, err := scheduler.Cron("0 6,12,18,0 * * *").Do(ComputeMarketRate)
	if err != nil {
		logger.Errorf("failed to schedule compute market rate task => %v\n", err)
	}

	// Refresh provision bucket priority queues every X hours
	_, err = scheduler.Cron(fmt.Sprintf("0 */%d * * *", conf.BucketQueueRebuildInterval)).
		Do(services.NewPriorityQueueService().ProcessBucketQueues)
	if err != nil {
		logger.Errorf("failed to schedule refresh priority queues task => %v\n", err)
	}

	// Retry failed webhook notifications every 1 minute
	_, err = scheduler.Cron("*/1 * * * *").Do(RetryFailedWebhookNotifications)
	if err != nil {
		logger.Errorf("cron.RetryFailedWebhookNotifications: %v\n", err)
	}

	// TODO: Run a task every 30 minutes for any failed order processing

	// Start scheduler
	scheduler.StartAsync()
}
