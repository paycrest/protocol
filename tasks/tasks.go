package tasks

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	"github.com/paycrest/protocol/services/contracts"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/logger"
	"github.com/shopspring/decimal"
)

var orderConf = config.OrderConfig()

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

// RetryStaleUserOperations retries stale user operations
func RetryStaleUserOperations() error {
	ctx := context.Background()
	orderService := services.NewOrderService()

	// Process initiated orders
	orders, err := storage.GetClient().PaymentOrder.
		Query().
		Where(func(s *sql.Selector) {
			ra := sql.Table(receiveaddress.Table)
			s.LeftJoin(ra).On(s.C(paymentorder.FieldReceiveAddressText), ra.C(receiveaddress.FieldAddress)).
				Where(sql.And(
					sql.EQ(s.C(paymentorder.FieldStatus), paymentorder.StatusInitiated),
					sql.EQ(ra.C(receiveaddress.FieldStatus), receiveaddress.StatusUsed),
					sql.IsNull(s.C(paymentorder.FieldGatewayID)),
				))
		}).
		All(ctx)
	if err != nil {
		return err
	}

	go func() {
		for _, order := range orders {
			orderAmountWithFees := order.Amount.Add(order.NetworkFee).Add(order.SenderFee).Add(order.ProtocolFee)
			if order.AmountPaid.GreaterThanOrEqual(orderAmountWithFees) {
				err := orderService.CreateOrder(ctx, order.ID)
				if err != nil {
					logger.Errorf("process task to create orders => %v", err)
				}
			}
		}
	}()

	// Revert order process
	orders, err = storage.GetClient().PaymentOrder.
		Query().
		Where(
			paymentorder.Or(
				paymentorder.StatusEQ(paymentorder.StatusInitiated),
				paymentorder.StatusEQ(paymentorder.StatusExpired),
			),
			paymentorder.AmountPaidGT(decimal.Zero),
			paymentorder.UpdatedAtLT(time.Now().Add(-5*time.Minute)),
		).
		WithReceiveAddress().
		All(ctx)
	if err != nil {
		return err
	}

	go func() {
		for _, order := range orders {
			if order.Edges.ReceiveAddress.Status == receiveaddress.StatusExpired || order.Edges.ReceiveAddress.Status == receiveaddress.StatusUsed {
				err := orderService.RevertOrder(ctx, order)
				if err != nil {
					logger.Errorf("process task to revert orders => %v", err)
				}
			}
		}
	}()

	// Settle order process
	lockOrders, err := storage.GetClient().LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusValidated),
			lockpaymentorder.HasFulfillmentWith(
				lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusSuccess),
			),
			lockpaymentorder.UpdatedAtLT(time.Now().Add(-5*time.Minute)),
		).
		All(ctx)
	if err != nil {
		return err
	}

	go func() {
		for _, order := range lockOrders {
			err := orderService.SettleOrder(ctx, order.ID)
			if err != nil {
				logger.Errorf("process order settlements task => %v", err)
			}
		}
	}()

	// Refund order process
	lockOrders, err = storage.GetClient().LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusPending),
			lockpaymentorder.CreatedAtLTE(time.Now().Add(-30*time.Minute)),
		).
		All(ctx)
	if err != nil {
		return err
	}

	go func() {
		for _, order := range lockOrders {
			err := orderService.RefundOrder(ctx, order.GatewayID)
			if err != nil {
				logger.Errorf("process order refunds task => %v", err)
			}
		}
	}()

	return nil
}

// IndexMissedBlocks indexes missed blocks
func IndexMissedBlocks() error {
	ctx := context.Background()
	orderService := services.NewOrderService()
	indexerService := services.NewIndexerService(orderService)

	networks, err := storage.Client.Network.Query().All(ctx)
	if err != nil {
		return fmt.Errorf("IndexMissedBlocks: %w", err)
	}

	var client types.RPCClient

	for _, network := range networks {
		// Index missed OrderCreated events
		orders, err := storage.GetClient().PaymentOrder.
			Query().
			Where(func(s *sql.Selector) {
				s.Where(sql.And(
					sql.EQ(s.C(paymentorder.FieldStatus), paymentorder.StatusPending),
					sql.IsNull(s.C(paymentorder.FieldGatewayID)),
					sql.LT(s.C(paymentorder.FieldUpdatedAt), time.Now().Add(-5*time.Minute))),
				)
			}).
			Order(ent.Asc(paymentorder.FieldBlockNumber)).
			All(ctx)
		if err != nil {
			logger.Errorf("IndexMissedBlocks: %v", err)
			continue
		}

		if len(orders) > 0 {
			retryErr := utils.Retry(3, 5*time.Second, func() error {
				client, err = types.NewEthClient(network.RPCEndpoint)
				return err
			})
			if retryErr != nil {
				logger.Errorf("IndexMissedBlocks: %v", err)
				continue
			}

			// Initialize contract filterer
			filterer, err := contracts.NewGatewayFilterer(orderConf.GatewayContractAddress, client)
			if err != nil {
				logger.Errorf("IndexMissedBlocks.NewGatewayFilterer: %v", err)
				return err
			}

			// Filter logs from the oldest indexed to the latest in the database
			toBlock := uint64(orders[len(orders)-1].BlockNumber)

			// Fetch logs
			var iter *contracts.GatewayOrderCreatedIterator
			retryErr = utils.Retry(3, 5*time.Second, func() error {
				var err error
				iter, err = filterer.FilterOrderCreated(&bind.FilterOpts{
					Start: uint64(orders[0].BlockNumber),
					End:   &toBlock,
				}, nil, nil, nil)
				return err
			})
			if retryErr != nil {
				logger.Errorf("IndexMissedBlocks.FilterOrderCreated: %v", retryErr)
				continue
			}

			// Iterate over logs
			for iter.Next() {
				err := indexerService.CreateLockPaymentOrder(ctx, client, network, iter.Event)
				if err != nil {
					logger.Errorf("IndexMissedBlocks.createOrder: %v", err)
					continue
				}
			}
		}
	}

	return nil
}

// HandleReceiveAddressValidity handles receive address validity
func HandleReceiveAddressValidity() error {
	ctx := context.Background()

	// Fetch expired receive addresses that are due for validity check
	addresses, err := storage.Client.ReceiveAddress.
		Query().
		Where(
			receiveaddress.ValidUntilLTE(time.Now()),
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
	// Fetch stable coin rate from third-party API Quidax (USDT)
	resp, err := utils.MakeJSONRequest(
		ctx,
		"GET",
		fmt.Sprintf("https://www.quidax.com/api/v1/markets/tickers/usdt%s", strings.ToLower(currency)),
		nil,
		nil,
	)
	if err != nil {
		return decimal.Zero, fmt.Errorf("ComputeMarketRate: %w", err)
	}

	price, err := decimal.NewFromString(resp["data"].(map[string]interface{})["ticker"].(map[string]interface{})["buy"].(string))
	if err != nil {
		return decimal.Zero, fmt.Errorf("ComputeMarketRate: %w", err)
	}

	return price, nil
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
		return fmt.Errorf("ComputeMarketRate: %w", err)
	}

	for _, currency := range currencies {
		// Fetch external rate
		externalRate, err := fetchExternalRate(ctx, currency.Code)
		if err != nil {
			return fmt.Errorf("ComputeMarketRate: %w", err)
		}

		// Fetch rates from token configs with fixed conversion rate
		token := "USDT"
		if config.ServerConfig().Environment != "production" {
			token = "6TEST"
		}
		tokenConfigs, err := storage.Client.ProviderOrderToken.
			Query().
			Where(
				providerordertoken.SymbolEQ(token),
				providerordertoken.ConversionRateTypeEQ(providerordertoken.ConversionRateTypeFixed),
			).
			All(ctx)
		if err != nil {
			return fmt.Errorf("ComputeMarketRate: %w", err)
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
			return fmt.Errorf("ComputeMarketRate: %w", err)
		}
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
		return fmt.Errorf("RetryFailedWebhookNotifications: %w", err)
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
				return fmt.Errorf("RetryFailedWebhookNotifications: %w", err)
			}

			continue
		}

		// Webhook notification was successful
		_, err := attempt.Update().
			SetStatus(webhookretryattempt.StatusSuccess).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("RetryFailedWebhookNotifications: %w", err)
		}
	}

	return nil
}

// StartCronJobs starts cron jobs
func StartCronJobs() {
	serverConf := config.ServerConfig()
	scheduler := gocron.NewScheduler(time.UTC)
	priorityQueue := services.NewPriorityQueueService()

	if serverConf.Environment != "production" {
		err := ComputeMarketRate()
		if err != nil {
			logger.Errorf("StartCronJobs: %v", err)
		}

		err = priorityQueue.ProcessBucketQueues()
		if err != nil {
			logger.Errorf("StartCronJobs: %v", err)
		}
	}

	// Compute market rate every 4 minutes
	_, err := scheduler.Cron("*/4 * * * *").Do(ComputeMarketRate)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Refresh provision bucket priority queues every X minutes
	_, err = scheduler.Cron(fmt.Sprintf("*/%d * * * *", orderConf.BucketQueueRebuildInterval)).
		Do(priorityQueue.ProcessBucketQueues)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Retry failed webhook notifications every 1 minute
	_, err = scheduler.Cron("*/1 * * * *").Do(RetryFailedWebhookNotifications)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Reassign pending order requests every 13 minutes
	_, err = scheduler.Cron("*/13 * * * *").Do(priorityQueue.ReassignPendingOrders)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Reassign unvalidated order requests every 21 minutes
	_, err = scheduler.Cron("*/21 * * * *").Do(priorityQueue.ReassignUnvalidatedLockOrders)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Handle receive address validity every 31 minutes
	_, err = scheduler.Cron("*/31 * * * *").Do(HandleReceiveAddressValidity)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Retry stale user operations every 5 minutes
	_, err = scheduler.Cron("*/5 * * * *").Do(RetryStaleUserOperations)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Index missed blocks every 7 minutes
	_, err = scheduler.Cron("*/7 * * * *").Do(IndexMissedBlocks)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Start scheduler
	scheduler.StartAsync()
}
