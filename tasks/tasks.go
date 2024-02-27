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
	"github.com/paycrest/protocol/ent/lockorderfulfillment"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	networkent "github.com/paycrest/protocol/ent/network"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/providerordertoken"
	"github.com/paycrest/protocol/ent/receiveaddress"
	"github.com/paycrest/protocol/ent/token"
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
	orderService := services.NewOrderService()
	indexerService := services.NewIndexerService(orderService)

	networks, err := storage.GetClient().Network.
		Query().
		Where(
			networkent.IsTestnetEQ(config.ServerConfig().Environment != "production"),
		).
		All(ctx)
	if err != nil {
		return err
	}

	for _, network := range networks {
		// Start listening for ERC20 transfer events. check for receive addresses updated within the last 12 hours
		twelveHoursAgo := time.Now().Add(-12 * time.Hour)

		addresses, err := storage.Client.ReceiveAddress.
			Query().
			Where(
				receiveaddress.HasPaymentOrderWith(
					paymentorder.HasTokenWith(
						token.HasNetworkWith(
							networkent.IDEQ(network.ID),
						),
					),
				),
				receiveaddress.UpdatedAtGT(twelveHoursAgo),
				receiveaddress.StatusNEQ(receiveaddress.StatusExpired),
			).
			Order(ent.Desc(receiveaddress.FieldLastIndexedBlock)).
			All(ctx)
		if err != nil {
			return err
		}

		for _, receiveAddress := range addresses {
			receiveAddress := receiveAddress

			go func() {
				_ = indexerService.IndexERC20Transfer(ctx, nil, receiveAddress)
			}()
		}

		// Start listening for order creation events
		go func(network *ent.Network) {
			err := indexerService.IndexOrderCreated(ctx, nil, network)
			if err != nil {
				logger.Errorf("process order deposits task => %v", err)
			}
		}(network)

		// Start listening for order settlement events
		go func(network *ent.Network) {
			err = indexerService.IndexOrderSettled(ctx, nil, network)
			if err != nil {
				logger.Errorf("process order settlements task => %v", err)
			}
		}(network)

		// Start listening for order refund events
		go func(network *ent.Network) {
			err = indexerService.IndexOrderRefunded(ctx, nil, network)
			if err != nil {
				logger.Errorf("process order refunds task => %v", err)
			}
		}(network)
	}

	return nil
}

// ProcessOrders processes orders
func ProcessOrders() error {
	ctx := context.Background()
	orderService := services.NewOrderService()

	// Create order process
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
			if order.AmountPaid.Equal(order.Amount) {
				err := orderService.CreateOrder(ctx, order.ID)
				if err != nil {
					logger.Errorf("process task to create orders => %v", err)
				}
			}
		}
	}()

	// // Revert order process
	// orders, err = storage.GetClient().PaymentOrder.
	// 	Query().
	// 	Where(
	// 		paymentorder.Or(
	// 			paymentorder.StatusEQ(paymentorder.StatusInitiated),
	// 			paymentorder.StatusEQ(paymentorder.StatusExpired),
	// 		),
	// 		paymentorder.AmountPaidGT(decimal.Zero),
	// 	).
	// 	All(ctx)
	// if err != nil {
	// 	return err
	// }

	// go func() {
	// 	for _, order := range orders {
	// 		fees := order.NetworkFee.Add(order.SenderFee)
	// 		orderAmountWithFees := order.Amount.Add(fees)
	// 		if !order.AmountPaid.Equal(orderAmountWithFees) && order.AmountReturned.Equal(decimal.Zero) {
	// 			err := orderService.RevertOrder(ctx, order, common.HexToAddress(order.FromAddress))
	// 			if err != nil {
	// 				logger.Errorf("process task to revert orders => %v", err)
	// 			}
	// 		}
	// 	}
	// }()

	// Settle order process
	lockOrders, err := storage.GetClient().LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusValidated),
			lockpaymentorder.HasFulfillmentWith(
				lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusSuccess),
			),
		).
		All(ctx)
	if err != nil {
		return err
	}

	go func() {
		for _, order := range lockOrders {
			orderService := services.NewOrderService()

			err := orderService.SettleOrder(ctx, order.ID)
			if err != nil {
				logger.Errorf("process order settlements task => %v", err)
			}
		}
	}()

	return nil
}

// ProcessOrderRefunds processes order refunds
func ProcessOrderRefunds() error {
	ctx := context.Background()

	// Refund orders
	orders, err := storage.GetClient().LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusPending),
			lockpaymentorder.CreatedAtLTE(time.Now().Add(-24*time.Hour)),
		).
		All(ctx)
	if err != nil {
		return err
	}

	go func() {
		for _, order := range orders {
			orderService := services.NewOrderService()

			err := orderService.RefundOrder(ctx, order.OrderID)
			if err != nil {
				logger.Errorf("process order refunds task => %v", err)
			}
		}
	}()

	return nil
}

// HandleReceiveAddressValidity handles receive address validity
func HandleReceiveAddressValidity() error {
	ctx := context.Background()
	orderConf := config.OrderConfig()

	// Fetch expired receive addresses that are due for validity check
	addresses, err := storage.Client.ReceiveAddress.
		Query().
		Where(
			receiveaddress.ValidUntilLTE(time.Now().Add(-orderConf.ReceiveAddressValidity)),
			receiveaddress.Or(
				receiveaddress.StatusNEQ(receiveaddress.StatusUsed),
				receiveaddress.And(
					receiveaddress.StatusEQ(receiveaddress.StatusUsed),
					receiveaddress.HasPaymentOrderWith(
						paymentorder.StatusEQ(paymentorder.StatusInitiated),
					),
				),
			),
		).
		WithPaymentOrder().
		All(ctx)
	if err != nil {
		return err
	}

	orderService := services.NewOrderService()
	indexerService := services.NewIndexerService(orderService)

	for _, address := range addresses {
		err := indexerService.HandleReceiveAddressValidity(ctx, address, address.Edges.PaymentOrder)
		if err != nil {
			continue
		}
	}

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

// fetchExternalRate fetches the external rate for a fiat currency
func fetchExternalRate(ctx context.Context, currency string) (decimal.Decimal, error) {
	// Fetch stable coin rate from third-party API Binance (USDT)
	resp, err := utils.MakeJSONRequest(
		ctx,
		"POST",
		"https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search",
		map[string]interface{}{
			"asset":     "USDT",
			"fiat":      currency,
			"tradeType": "SELL",
			"page":      1,
			"rows":      20,
		},
		nil,
	)
	if err != nil {
		logger.Errorf("ComputeMarketRate: %v", err)
		return decimal.Zero, err
	}

	// Access the data array
	data, ok := resp["data"].([]interface{})
	if !ok || len(data) == 0 {
		return decimal.Zero, fmt.Errorf("ComputeMarketRate: No data in the response")
	}

	// Loop through the data array and extract prices
	var prices []decimal.Decimal
	for _, item := range data {
		adv, ok := item.(map[string]interface{})["adv"].(map[string]interface{})
		if !ok {
			logger.Errorf("ComputeMarketRate: adv not found or not a map.")
			continue
		}

		price, err := decimal.NewFromString(adv["price"].(string))
		if err != nil {
			logger.Errorf("ComputeMarketRate: %v", err)
			continue
		}

		prices = append(prices, price)
	}

	// Calculate and return the median
	return utils.Median(prices), nil
}

// ComputeMarketRate computes the market price for fiat currencies
func ComputeMarketRate() error {
	ctx := context.Background()
	orderConf := config.OrderConfig()

	// Fetch all fiat currencies
	currencies, err := storage.Client.FiatCurrency.
		Query().
		Where(fiatcurrency.IsEnabledEQ(true)).
		All(ctx)
	if err != nil {
		return err
	}

	for _, currency := range currencies {
		// Fetch external rate
		externalRate, err := fetchExternalRate(ctx, currency.Code)
		if err != nil {
			logger.Errorf("ComputeMarketRate: %v", err)
			return err
		}

		// Fetch rates from token configs with fixed conversion rate
		tokenConfigs, err := storage.Client.ProviderOrderToken.
			Query().
			Where(
				providerordertoken.SymbolEQ("USDT"),
				providerordertoken.ConversionRateTypeEQ(providerordertoken.ConversionRateTypeFixed),
			).
			All(ctx)
		if err != nil {
			logger.Errorf("ComputeMarketRate: %v", err)
		}

		var rates []decimal.Decimal
		for _, tokenConfig := range tokenConfigs {
			rates = append(rates, tokenConfig.FixedConversionRate)
		}

		// Calculate median
		median := utils.Median(rates)

		// Check the median rate against the external rate to ensure it's not too far off
		percentDeviation := utils.AbsPercentageDeviation(externalRate, median)
		if percentDeviation.GreaterThan(orderConf.PercentDeviationFromExternalRate) {
			median = externalRate
		}

		// Update currency with median rate
		_, err = storage.Client.FiatCurrency.
			UpdateOneID(currency.ID).
			SetMarketRate(median).
			Save(ctx)
		if err != nil {
			logger.Errorf("ComputeMarketRate: %v", err)
			return err
		}

		logger.Infof("Computed market rate for %s: %s\n", currency.Code, median.String())
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
		logger.Errorf("RetryFailedWebhookNotifications: %v", err)
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
				logger.Errorf("RetryFailedWebhookNotifications: %v", err)
			}

			continue
		}

		// Webhook notification was successful
		_, err := attempt.Update().
			SetStatus(webhookretryattempt.StatusSuccess).
			Save(ctx)
		if err != nil {
			logger.Errorf("RetryFailedWebhookNotifications: %v", err)
		}
	}

	return nil
}

// StartCronJobs starts cron jobs
func StartCronJobs() {
	orderConf := config.OrderConfig()
	serverConf := config.ServerConfig()
	scheduler := gocron.NewScheduler(time.UTC)
	priorityQueue := services.NewPriorityQueueService()

	if serverConf.Environment != "production" {
		err := ComputeMarketRate()
		if err != nil {
			logger.Errorf("failed to compute market rate => %v", err)
		}

		err = priorityQueue.ProcessBucketQueues()
		if err != nil {
			logger.Errorf("failed to process bucket queues => %v", err)
		}
	}

	// Compute market rate every 10 minutes
	_, err := scheduler.Cron("*/10 * * * *").Do(ComputeMarketRate)
	if err != nil {
		logger.Errorf("failed to schedule compute market rate task => %v", err)
	}

	// Refresh provision bucket priority queues every X minutes
	_, err = scheduler.Cron(fmt.Sprintf("*/%d * * * *", orderConf.BucketQueueRebuildInterval)).
		Do(priorityQueue.ProcessBucketQueues)
	if err != nil {
		logger.Errorf("failed to schedule refresh priority queues task => %v", err)
	}

	// Retry failed webhook notifications every 1 minute
	_, err = scheduler.Cron("*/1 * * * *").Do(RetryFailedWebhookNotifications)
	if err != nil {
		logger.Errorf("cron.RetryFailedWebhookNotifications: %v", err)
	}

	// Reassign pending order requests every 15 minutes
	_, err = scheduler.Cron("*/15 * * * *").Do(priorityQueue.ReassignPendingOrders)
	if err != nil {
		logger.Errorf("cron.ReassignPendingOrders: %v", err)
	}

	// Reassign unvalidated order requests every 20 minutes
	_, err = scheduler.Cron("*/20 * * * *").Do(priorityQueue.ReassignUnvalidatedLockOrders)
	if err != nil {
		logger.Errorf("cron.ReassignUnvalidatedLockOrders: %v", err)
	}

	// Handle receive address validity every 30 minutes
	_, err = scheduler.Cron("*/30 * * * *").Do(HandleReceiveAddressValidity)
	if err != nil {
		logger.Errorf("cron.HandleReceiveAddressValidity: %v", err)
	}

	// Process order refunds once a day
	_, err = scheduler.Cron("0 0 * * *").Do(ProcessOrderRefunds)
	if err != nil {
		logger.Errorf("cron.ProcessOrderRefunds: %v", err)
	}

	// Start scheduler
	scheduler.StartAsync()
}
