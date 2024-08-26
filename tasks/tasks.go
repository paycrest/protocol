package tasks

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/ent/lockorderfulfillment"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	networkent "github.com/paycrest/protocol/ent/network"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/paymentorderrecipient"
	"github.com/paycrest/protocol/ent/providerordertoken"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/receiveaddress"
	"github.com/paycrest/protocol/ent/senderprofile"
	"github.com/paycrest/protocol/ent/token"
	"github.com/paycrest/protocol/ent/transactionlog"
	"github.com/paycrest/protocol/ent/webhookretryattempt"
	"github.com/paycrest/protocol/services"
	orderService "github.com/paycrest/protocol/services/order"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/logger"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

var orderConf = config.OrderConfig()
var serverConf = config.ServerConfig()

var rpcClients = map[string]types.RPCClient{}

// setRPCClients connects to the RPC endpoints of all networks
func setRPCClients(ctx context.Context) ([]*ent.Network, error) {
	isTestnet := false
	if serverConf.Environment != "production" {
		isTestnet = true
	}

	networks, err := storage.Client.Network.
		Query().
		Where(networkent.IsTestnetEQ(isTestnet)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("setRPCClients.fetchNetworks: %w", err)
	}

	// Connect to RPC endpoint
	var client types.RPCClient
	for _, network := range networks {
		if rpcClients[network.Identifier] == nil && !strings.HasPrefix(network.Identifier, "tron") {
			retryErr := utils.Retry(3, 1*time.Second, func() error {
				client, err = types.NewEthClient(network.RPCEndpoint)
				return err
			})
			if retryErr != nil {
				logger.Errorf("failed to connect to %s RPC %v", network.Identifier, retryErr)
				continue
			}

			rpcClients[network.Identifier] = client
		}
	}

	return networks, nil
}

// RetryStaleUserOperations retries stale user operations
func RetryStaleUserOperations() error {
	ctx := context.Background()
	var wg sync.WaitGroup

	// Establish RPC connections
	_, err := setRPCClients(ctx)
	if err != nil {
		return fmt.Errorf("RetryStaleUserOperations: %w", err)
	}

	// Process initiated orders
	orders, err := storage.Client.PaymentOrder.
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
		Where(
			paymentorder.Or(
				paymentorder.UpdatedAtGTE(time.Now().Add(-10*time.Minute)),
				paymentorder.HasRecipientWith(
					paymentorderrecipient.MemoHasPrefix("P#P"),
				),
			)).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		All(ctx)
	if err != nil {
		return fmt.Errorf("RetryStaleUserOperations: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, order := range orders {
			orderAmountWithFees := order.Amount.Add(order.NetworkFee).Add(order.SenderFee).Add(order.ProtocolFee)
			if order.AmountPaid.GreaterThanOrEqual(orderAmountWithFees) {
				var service types.OrderService
				if strings.HasPrefix(order.Edges.Token.Edges.Network.Identifier, "tron") {
					service = orderService.NewOrderTron()
				} else {
					service = orderService.NewOrderEVM()
				}
				err := service.CreateOrder(ctx, rpcClients[order.Edges.Token.Edges.Network.Identifier], order.ID)
				if err != nil {
					logger.Errorf("RetryStaleUserOperations.CreateOrder %v", err)
				}
			}
		}
	}()

	// Revert order process
	orders, err = storage.Client.PaymentOrder.
		Query().
		Where(
			paymentorder.Or(
				paymentorder.StatusEQ(paymentorder.StatusInitiated),
				paymentorder.StatusEQ(paymentorder.StatusExpired),
			),
			paymentorder.AmountPaidGT(decimal.Zero),
			paymentorder.UpdatedAtLT(time.Now().Add(-10*time.Minute)),
		).
		WithReceiveAddress().
		WithRecipient().
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		All(ctx)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, order := range orders {
			if order.Edges.ReceiveAddress.Status == receiveaddress.StatusExpired || order.Edges.ReceiveAddress.Status == receiveaddress.StatusUsed || order.Edges.ReceiveAddress.Status == receiveaddress.StatusPartial {
				var service types.OrderService
				if strings.HasPrefix(order.Edges.Token.Edges.Network.Identifier, "tron") {
					service = orderService.NewOrderTron()
				} else {
					service = orderService.NewOrderEVM()
				}
				err := service.RevertOrder(ctx, rpcClients[order.Edges.Token.Edges.Network.Identifier], order)
				if err != nil {
					logger.Errorf("RetryStaleUserOperations.RevertOrder: %v", err)
				}
			}
		}
	}()

	// Settle order process
	lockOrders, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusValidated),
			lockpaymentorder.HasFulfillmentsWith(
				lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusSuccess),
			),
			lockpaymentorder.UpdatedAtLT(time.Now().Add(-5*time.Minute)),
		).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		All(ctx)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, order := range lockOrders {
			var service types.OrderService
			if strings.HasPrefix(order.Edges.Token.Edges.Network.Identifier, "tron") {
				service = orderService.NewOrderTron()
			} else {
				service = orderService.NewOrderEVM()
			}
			err := service.SettleOrder(ctx, rpcClients[order.Edges.Token.Edges.Network.Identifier], order.ID)
			if err != nil {
				logger.Errorf("RetryStaleUserOperations.SettleOrder: %v", err)
			}
		}
	}()

	// Refund order process
	lockOrders, err = storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.GatewayIDNEQ(""),
			lockpaymentorder.Or(
				lockpaymentorder.And(
					lockpaymentorder.Or(
						lockpaymentorder.StatusEQ(lockpaymentorder.StatusPending),
						lockpaymentorder.StatusEQ(lockpaymentorder.StatusCancelled),
					),
					lockpaymentorder.CreatedAtLTE(time.Now().Add(-30*time.Minute)),
				),
				lockpaymentorder.And(
					lockpaymentorder.HasProviderWith(
						providerprofile.VisibilityModeEQ(providerprofile.VisibilityModePrivate),
					),
					lockpaymentorder.StatusEQ(lockpaymentorder.StatusFulfilled),
					lockpaymentorder.HasFulfillmentsWith(
						lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusFailed),
						lockorderfulfillment.Not(lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusSuccess)),
						lockorderfulfillment.Not(lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusPending)),
					),
				),
			),
		).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		All(ctx)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, order := range lockOrders {
			var service types.OrderService
			if strings.HasPrefix(order.Edges.Token.Edges.Network.Identifier, "tron") {
				service = orderService.NewOrderTron()
			} else {
				service = orderService.NewOrderEVM()
			}
			err := service.RefundOrder(ctx, rpcClients[order.Edges.Token.Edges.Network.Identifier], order.GatewayID)
			if err != nil {
				logger.Errorf("RetryStaleUserOperations.RefundOrder: %v", err)
			}
		}
	}()

	return nil
}

// IndexBlockchainEvents indexes missed blocks
func IndexBlockchainEvents() error {
	ctx := context.Background()
	var wg sync.WaitGroup

	time.Sleep(100 * time.Millisecond) // to keep out of sync with other tasks

	// Establish RPC connections
	networks, err := setRPCClients(ctx)
	if err != nil {
		return fmt.Errorf("IndexBlockchainEvents: %w", err)
	}

	// Index ERC20 transfer events
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = utils.Retry(3, 2*time.Second, func() error {
			orders, err := storage.Client.PaymentOrder.
				Query().
				Where(
					paymentorder.StatusEQ(paymentorder.StatusInitiated),
					paymentorder.HasReceiveAddressWith(
						receiveaddress.Or(
							receiveaddress.StatusEQ(receiveaddress.StatusUnused),
							receiveaddress.StatusEQ(receiveaddress.StatusPartial),
						),
						receiveaddress.ValidUntilGT(time.Now()),
					),
				).
				WithToken(func(tq *ent.TokenQuery) {
					tq.WithNetwork()
				}).
				WithReceiveAddress().
				WithRecipient().
				All(ctx)
			if err != nil {
				logger.Errorf("IndexBlockchainEvents: %v", err)
			}

			if len(orders) > 0 {
				for _, order := range orders {
					if strings.HasPrefix(order.Edges.Token.Edges.Network.Identifier, "tron") {
						indexerService := services.NewIndexerService(orderService.NewOrderTron())
						err := indexerService.IndexTRC20Transfer(ctx, order)
						if err != nil {
							continue
						}
					} else {
						indexerService := services.NewIndexerService(orderService.NewOrderEVM())
						err := indexerService.IndexERC20Transfer(ctx, rpcClients[order.Edges.Token.Edges.Network.Identifier], order)
						if err != nil {
							continue
						}
					}
				}
			}

			return fmt.Errorf("trigger retry")
		})
	}()

	// Index OrderCreated events
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(500 * time.Millisecond)
		_ = utils.Retry(3, 2*time.Second, func() error {
			for _, network := range networks {
				// Index events triggered from API
				orders, err := storage.Client.PaymentOrder.
					Query().
					Where(func(s *sql.Selector) {
						lpo := sql.Table(lockpaymentorder.Table)
						s.Where(sql.And(
							sql.EQ(s.C(paymentorder.FieldStatus), paymentorder.StatusPending),
							sql.Or(
								sql.NotExists(
									sql.Select().
										From(lpo).
										Where(sql.ColumnsEQ(s.C(paymentorder.FieldGatewayID), lpo.C(lockpaymentorder.FieldGatewayID))),
								),
								sql.IsNull(s.C(paymentorder.FieldGatewayID)),
							),
							sql.GT(s.C(paymentorder.FieldBlockNumber), 0),
						))
					}).
					Where(paymentorder.HasTokenWith(token.HasNetworkWith(networkent.IDEQ(network.ID)))).
					WithReceiveAddress().
					WithToken(func(tq *ent.TokenQuery) {
						tq.WithNetwork()
					}).
					Order(ent.Asc(paymentorder.FieldBlockNumber)).
					All(ctx)
				if err != nil {
					continue
				}

				if len(orders) > 0 {
					for _, order := range orders {
						if strings.HasPrefix(network.Identifier, "tron") {
							indexerService := services.NewIndexerService(orderService.NewOrderTron())
							err := indexerService.IndexOrderCreatedTron(ctx, order)
							if err != nil {
								continue
							}
						} else {
							indexerService := services.NewIndexerService(orderService.NewOrderEVM())
							err := indexerService.IndexOrderCreated(ctx, rpcClients[network.Identifier], network, order.Edges.ReceiveAddress.Address)
							if err != nil {
								continue
							}
						}
					}
				}

				// Index events triggered from Gateway contract
				if !strings.HasPrefix(network.Identifier, "tron") {
					indexerService := services.NewIndexerService(orderService.NewOrderEVM())
					err = indexerService.IndexOrderCreated(ctx, rpcClients[network.Identifier], network, "")
					if err != nil {
						continue
					}
				}
			}

			return fmt.Errorf("trigger retry")
		})
	}()

	// Index OrderSettled events
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(1000 * time.Millisecond)
		_ = utils.Retry(3, 2*time.Second, func() error {
			lockOrders, err := storage.Client.LockPaymentOrder.
				Query().
				Where(func(s *sql.Selector) {
					po := sql.Table(paymentorder.Table)
					s.LeftJoin(po).On(s.C(lockpaymentorder.FieldGatewayID), po.C(paymentorder.FieldGatewayID)).
						Where(sql.Or(
							sql.EQ(s.C(lockpaymentorder.FieldStatus), lockpaymentorder.StatusValidated),
							sql.And(
								sql.EQ(po.C(paymentorder.FieldStatus), paymentorder.StatusPending),
								sql.EQ(s.C(lockpaymentorder.FieldStatus), lockpaymentorder.StatusSettled),
							)),
						)
				}).
				WithToken(func(tq *ent.TokenQuery) {
					tq.WithNetwork()
				}).
				Order(ent.Asc(lockpaymentorder.FieldBlockNumber)).
				All(ctx)
			if err != nil {
				logger.Errorf("IndexBlockchainEvents: %v", err)
			}

			if len(lockOrders) > 0 {
				for _, order := range lockOrders {
					if strings.HasPrefix(order.Edges.Token.Edges.Network.Identifier, "tron") {
						indexerService := services.NewIndexerService(orderService.NewOrderTron())
						err := indexerService.IndexOrderSettledTron(ctx, order)
						if err != nil {
							continue
						}
					} else {
						indexerService := services.NewIndexerService(orderService.NewOrderEVM())
						err := indexerService.IndexOrderSettled(ctx, rpcClients[order.Edges.Token.Edges.Network.Identifier], order.Edges.Token.Edges.Network, order.GatewayID)
						if err != nil {
							continue
						}
					}
				}
			}

			return fmt.Errorf("trigger retry")
		})
	}()

	// Index OrderRefunded events
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(1500 * time.Millisecond)
		_ = utils.Retry(3, 2*time.Second, func() error {
			lockOrders, err := storage.Client.LockPaymentOrder.
				Query().
				Where(func(s *sql.Selector) {
					po := sql.Table(paymentorder.Table)
					s.LeftJoin(po).On(s.C(lockpaymentorder.FieldGatewayID), po.C(paymentorder.FieldGatewayID)).
						Where(sql.Or(
							sql.And(
								sql.EQ(s.C(lockpaymentorder.FieldStatus), lockpaymentorder.StatusPending),
								sql.LT(s.C(lockpaymentorder.FieldCreatedAt), time.Now().Add(-35*time.Minute)),
							),
							sql.And(
								sql.EQ(po.C(paymentorder.FieldStatus), paymentorder.StatusPending),
								sql.EQ(s.C(lockpaymentorder.FieldStatus), lockpaymentorder.StatusRefunded),
							),
						))
				}).
				WithToken(func(tq *ent.TokenQuery) {
					tq.WithNetwork()
				}).
				Order(ent.Asc(lockpaymentorder.FieldBlockNumber)).
				All(ctx)
			if err != nil {
				logger.Errorf("IndexBlockchainEvents: %v", err)
			}

			if len(lockOrders) > 0 {
				for _, order := range lockOrders {
					if strings.HasPrefix(order.Edges.Token.Edges.Network.Identifier, "tron") {
						indexerService := services.NewIndexerService(orderService.NewOrderTron())
						err := indexerService.IndexOrderRefundedTron(ctx, order)
						if err != nil {
							continue
						}
					} else {
						indexerService := services.NewIndexerService(orderService.NewOrderEVM())
						err := indexerService.IndexOrderRefunded(ctx, rpcClients[order.Edges.Token.Edges.Network.Identifier], order.Edges.Token.Edges.Network, order.GatewayID)
						if err != nil {
							continue
						}
					}
				}
			}

			return fmt.Errorf("trigger retry")
		})
	}()

	return nil
}

// ReassignPendingOrders reassigns declined order requests to providers
func ReassignPendingOrders() {
	ctx := context.Background()

	// Query pending lock orders
	lockOrders, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusPending),
			lockpaymentorder.Not(lockpaymentorder.HasFulfillments()),
		).
		WithToken().
		WithProvider().
		WithProvisionBucket(
			func(pbq *ent.ProvisionBucketQuery) {
				pbq.WithCurrency()
			},
		).
		All(ctx)
	if err != nil {
		logger.Errorf("ReassignPendingOrders.db: %v", err)
		return
	}

	// Check if order_request_<order_id> exists in Redis
	for _, order := range lockOrders {
		orderKey := fmt.Sprintf("order_request_%s", order.ID)
		exists, err := storage.RedisClient.Exists(ctx, orderKey).Result()
		if err != nil {
			logger.Errorf("ReassignPendingOrders.redis: %v", err)
			continue
		}

		if exists == 0 {
			// Order request doesn't exist in Redis, reassign the order
			lockPaymentOrder := types.LockPaymentOrderFields{
				ID:                order.ID,
				Token:             order.Edges.Token,
				GatewayID:         order.GatewayID,
				Amount:            order.Amount,
				Rate:              order.Rate,
				BlockNumber:       order.BlockNumber,
				Institution:       order.Institution,
				AccountIdentifier: order.AccountIdentifier,
				AccountName:       order.AccountName,
				Memo:              order.Memo,
				ProvisionBucket:   order.Edges.ProvisionBucket,
			}

			if order.Edges.Provider != nil {
				lockPaymentOrder.ProviderID = order.Edges.Provider.ID
			}

			err := services.NewPriorityQueueService().AssignLockPaymentOrder(ctx, lockPaymentOrder)
			if err != nil {
				logger.Errorf("failed to reassign declined order request: %v", err)
			}
		}
	}
}

// ReassignUnfulfilledLockOrders reassigns lockOrder unfulfilled within a time frame.
func ReassignUnfulfilledLockOrders() {
	ctx := context.Background()

	// Query unfulfilled lock orders.
	lockOrders, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.Or(
				lockpaymentorder.Not(lockpaymentorder.HasFulfillments()),
				lockpaymentorder.HasFulfillmentsWith(
					lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusFailed),
					lockorderfulfillment.Not(lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusSuccess)),
					lockorderfulfillment.Not(lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusPending)),
				),
			),
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusProcessing),
			lockpaymentorder.Or(
				lockpaymentorder.UpdatedAtLTE(time.Now().Add(-orderConf.OrderFulfillmentValidity*time.Minute)),
				lockpaymentorder.HasFulfillmentsWith(
					lockorderfulfillment.CreatedAtLTE(time.Now().Add(-orderConf.OrderFulfillmentValidity*time.Minute)),
				),
			),
		).
		WithToken().
		WithProvider().
		WithProvisionBucket(func(pbq *ent.ProvisionBucketQuery) {
			pbq.WithCurrency()
		}).
		All(ctx)
	if err != nil {
		logger.Errorf("ReassignUnfulfilledLockOrders: %v", err)
		return
	}

	// Unassign unfulfilled lock orders.
	_, err = storage.Client.LockPaymentOrder.
		Update().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusProcessing),
			lockpaymentorder.UpdatedAtLTE(time.Now().Add(-orderConf.OrderFulfillmentValidity*time.Minute)),
			lockpaymentorder.Or(
				lockpaymentorder.Not(lockpaymentorder.HasFulfillments()),
				lockpaymentorder.HasFulfillmentsWith(
					lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusFailed),
					lockorderfulfillment.Not(lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusSuccess)),
					lockorderfulfillment.Not(lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusPending)),
				),
			),
		).
		SetStatus(lockpaymentorder.StatusPending).
		Save(ctx)
	if err != nil {
		logger.Errorf("ReassignUnfulfilledLockOrders: %v", err)
		return
	}

	for _, order := range lockOrders {
		lockPaymentOrder := types.LockPaymentOrderFields{
			ID:                order.ID,
			Token:             order.Edges.Token,
			GatewayID:         order.GatewayID,
			Amount:            order.Amount,
			Rate:              order.Rate,
			BlockNumber:       order.BlockNumber,
			Institution:       order.Institution,
			AccountIdentifier: order.AccountIdentifier,
			AccountName:       order.AccountName,
			ProviderID:        order.Edges.Provider.ID,
			Memo:              order.Memo,
			ProvisionBucket:   order.Edges.ProvisionBucket,
		}

		err := services.NewPriorityQueueService().AssignLockPaymentOrder(ctx, lockPaymentOrder)
		if err != nil {
			logger.Errorf("ReassignUnfulfilledLockOrders.AssignLockPaymentOrder: %s => %v", order.GatewayID, err)
		}
	}
}

// ReassignUnvalidatedLockOrders reassigns or refunds unvalidated lock orders to providers
func ReassignUnvalidatedLockOrders() {
	ctx := context.Background()

	// Query unvalidated lock orders.
	lockOrders, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusFulfilled),
			lockpaymentorder.Or(
				lockpaymentorder.HasFulfillmentsWith(
					lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusFailed),
					lockorderfulfillment.Not(lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusSuccess)),
					lockorderfulfillment.Not(lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusPending)),
				),
				lockpaymentorder.And(
					lockpaymentorder.HasFulfillmentsWith(
						lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusPending),
					),
					lockpaymentorder.HasFulfillmentsWith(
						lockorderfulfillment.UpdatedAtLTE(time.Now().Add(-orderConf.OrderFulfillmentValidity*time.Minute)),
						lockorderfulfillment.Not(lockorderfulfillment.UpdatedAtGT(time.Now().Add(-orderConf.OrderFulfillmentValidity*time.Minute))),
					),
				),
			),
		).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		WithProvider(
			func(pq *ent.ProviderProfileQuery) {
				pq.WithAPIKey()
			},
		).
		WithFulfillments().
		WithProvisionBucket(
			func(pbq *ent.ProvisionBucketQuery) {
				pbq.WithCurrency()
			},
		).
		All(ctx)
	if err != nil {
		logger.Errorf("ReassignUnvalidatedLockOrders.db: %v", err)
		return
	}

	for _, order := range lockOrders {
		for _, fulfillment := range order.Edges.Fulfillments {
			if fulfillment.ValidationStatus == lockorderfulfillment.ValidationStatusPending {
				// TODO: use auth
				// // Compute HMAC
				// decodedSecret, err := base64.StdEncoding.DecodeString(order.Edges.Provider.Edges.APIKey.Secret)
				// if err != nil {
				// 	logger.Errorf("ReassignUnvalidatedLockOrders: %v", err)
				// 	return
				// }
				// decryptedSecret, err := cryptoUtils.DecryptPlain(decodedSecret)
				// if err != nil {
				// 	logger.Errorf("ReassignUnvalidatedLockOrders: %v", err)
				// 	return
				// }

				// payload := map[string]interface{}{}

				// signature := tokenUtils.GenerateHMACSignature(payload, string(decryptedSecret))

				// Send GET request to the provider's node
				res, err := fastshot.NewClient(order.Edges.Provider.HostIdentifier).
					Config().SetTimeout(30 * time.Second).
					// Header().Add("X-Request-Signature", signature).
					Build().GET(fmt.Sprintf("/tx_status/%s/%s", fulfillment.Psp, fulfillment.TxID)).
					Send()
				if err != nil {
					logger.Errorf("ReassignUnvalidatedLockOrders: %v", err)
					continue
				}

				data, err := utils.ParseJSONResponse(res.RawResponse)
				if err != nil {
					logger.Errorf("ReassignUnvalidatedLockOrders: %v", err)
					continue
				}

				status := data["data"].(map[string]interface{})["status"].(string)

				if status == "failed" {
					_, err = storage.Client.LockOrderFulfillment.
						UpdateOneID(fulfillment.ID).
						SetValidationStatus(lockorderfulfillment.ValidationStatusFailed).
						SetValidationError(data["data"].(map[string]interface{})["error"].(string)).
						Save(ctx)
					if err != nil {
						logger.Errorf("ReassignUnvalidatedLockOrders.UpdateFulfillmentStatusFailed: %v", err)
						continue
					}

					_, err = order.Update().
						SetStatus(lockpaymentorder.StatusFulfilled).
						Save(ctx)
					if err != nil {
						logger.Errorf("ReassignUnvalidatedLockOrders.UpdateOrderStatusFulfilled: %v", err)
						continue
					}

				} else if status == "success" {
					_, err = storage.Client.LockOrderFulfillment.
						UpdateOneID(fulfillment.ID).
						SetValidationStatus(lockorderfulfillment.ValidationStatusSuccess).
						Save(ctx)
					if err != nil {
						logger.Errorf("ReassignUnvalidatedLockOrders.UpdateFulfillmentStatusSuccess: %v", err)
						continue
					}

					transactionLog, err := storage.Client.TransactionLog.Create().
						SetStatus(transactionlog.StatusOrderValidated).
						SetNetwork(fulfillment.Edges.Order.Edges.Token.Edges.Network.Identifier).
						SetMetadata(map[string]interface{}{
							"TransactionID": fulfillment.TxID,
							"PSP":           fulfillment.Psp,
						}).
						Save(ctx)
					if err != nil {
						logger.Errorf("ReassignUnvalidatedLockOrders.CreateTransactionLog: %v", err)
						continue
					}

					_, err = order.Update().
						SetStatus(lockpaymentorder.StatusValidated).
						AddTransactions(transactionLog).
						Save(ctx)
					if err != nil {
						logger.Errorf("ReassignUnvalidatedLockOrders.UpdateOrderStatusValidated: %v", err)
						continue
					}
				}

			} else {
				if order.Edges.Provider.VisibilityMode != providerprofile.VisibilityModePrivate {
					lockPaymentOrder := types.LockPaymentOrderFields{
						ID:                order.ID,
						Token:             order.Edges.Token,
						GatewayID:         order.GatewayID,
						Amount:            order.Amount,
						Rate:              order.Rate,
						BlockNumber:       order.BlockNumber,
						Institution:       order.Institution,
						AccountIdentifier: order.AccountIdentifier,
						AccountName:       order.AccountName,
						ProviderID:        "",
						Memo:              order.Memo,
						ProvisionBucket:   order.Edges.ProvisionBucket,
					}

					err := services.NewPriorityQueueService().AssignLockPaymentOrder(ctx, lockPaymentOrder)
					if err != nil {
						logger.Errorf("ReassignUnvalidatedLockOrders.AssignLockPaymentOrder: %v", err)
					}
				}
			}
		}
	}
}

// ReassignStaleOrderRequest reassigns expired order requests to providers
func ReassignStaleOrderRequest(ctx context.Context, orderRequestChan <-chan *redis.Message) {
	for msg := range orderRequestChan {
		key := strings.Split(msg.Payload, "_")
		orderID := key[len(key)-1]

		orderUUID, err := uuid.Parse(orderID)
		if err != nil {
			logger.Errorf("ReassignStaleOrderRequest: %v", err)
			continue
		}

		// Get the order from the database
		order, err := storage.Client.LockPaymentOrder.
			Query().
			Where(
				lockpaymentorder.IDEQ(orderUUID),
			).
			WithProvisionBucket().
			Only(ctx)
		if err != nil {
			logger.Errorf("ReassignStaleOrderRequest: %v", err)
			continue
		}

		orderFields := types.LockPaymentOrderFields{
			ID:                order.ID,
			GatewayID:         order.GatewayID,
			Amount:            order.Amount,
			Rate:              order.Rate,
			BlockNumber:       order.BlockNumber,
			Institution:       order.Institution,
			AccountIdentifier: order.AccountIdentifier,
			AccountName:       order.AccountName,
			Memo:              order.Memo,
			ProvisionBucket:   order.Edges.ProvisionBucket,
		}

		// Assign the order to a provider
		err = services.NewPriorityQueueService().AssignLockPaymentOrder(ctx, orderFields)
		if err != nil {
			logger.Errorf("ReassignStaleOrderRequest.AssignLockPaymentOrder: %v", err)
		}
	}
}

func FixDatabaseMisHap() error {
	ctx := context.Background()

	// Restore initiated orders
	// _, err := storage.Client.PaymentOrder.
	// 	Update().
	// 	Where(
	// 		paymentorder.AmountPaidEQ(decimal.Zero),
	// 	).
	// 	SetStatus(paymentorder.StatusInitiated).
	// 	SetBlockNumber(0).
	// 	SetPercentSettled(decimal.NewFromInt(0)).
	// 	SetTxHash("").
	// 	Save(ctx)
	// if err != nil {
	// 	logger.Errorf("FixDatabaseMisHap: %v", err)
	// }

	// Restore pending orders
	orders, err := storage.Client.PaymentOrder.
		Query().
		Where(
			paymentorder.Or(
				paymentorder.And(
					paymentorder.AmountPaidGT(decimal.Zero),
					paymentorder.StatusEQ(paymentorder.StatusPending),
				),
				// paymentorder.And(
				// 	paymentorder.AmountPaidGT(decimal.Zero),
				// 	paymentorder.StatusEQ(paymentorder.StatusSettled),
				// 	paymentorder.GatewayIDNEQ(""),
				// ),
			),
		).
		WithRecipient().
		All(ctx)
	if err != nil {
		logger.Errorf("FixDatabaseMisHap: %v", err)
	}

	for _, order := range orders {
		lockOrder, err := storage.Client.LockPaymentOrder.
			Query().
			Where(
				lockpaymentorder.MemoEQ(order.Edges.Recipient.Memo),
				lockpaymentorder.MemoHasPrefix("P#P"),
				lockpaymentorder.Or(
					lockpaymentorder.StatusEQ(lockpaymentorder.StatusSettled),
					lockpaymentorder.StatusEQ(lockpaymentorder.StatusRefunded),
				),
			).
			First(ctx)
		if err != nil {
			continue
		}
		if lockOrder != nil && order.AmountPaid.GreaterThanOrEqual(order.Amount.Add(order.NetworkFee).Add(order.SenderFee).Add(order.ProtocolFee)) {
			// var status paymentorder.Status
			// var blockNumber int64
			// var percentSettled decimal.Decimal
			// var txHash string

			if lockOrder.Status == lockpaymentorder.StatusSettled {
				// Fix settled orders without gatewayId
				_, err := order.Update().
					SetStatus(paymentorder.StatusSettled).
					SetBlockNumber(lockOrder.BlockNumber).
					SetPercentSettled(decimal.NewFromInt(100)).
					SetTxHash(lockOrder.TxHash).
					SetGatewayID(lockOrder.GatewayID).
					Save(ctx)
				if err != nil {
					logger.Errorf("FixDatabaseMisHap: %v", err)
					continue
				}
			} else if lockOrder.Status == lockpaymentorder.StatusRefunded {
				// Fix refunded orders without gatewayId
				_, err := order.Update().
					SetStatus(paymentorder.StatusRefunded).
					SetBlockNumber(lockOrder.BlockNumber).
					SetPercentSettled(decimal.NewFromInt(0)).
					SetTxHash(lockOrder.TxHash).
					SetGatewayID(lockOrder.GatewayID).
					Save(ctx)
				if err != nil {
					logger.Errorf("FixDatabaseMisHap: %v", err)
					continue
				}
			}

			// if lockOrder.Status == lockpaymentorder.StatusRefunded {
			// 	status = paymentorder.StatusRefunded
			// 	blockNumber = lockOrder.BlockNumber
			// 	percentSettled = decimal.NewFromInt(0)
			// 	txHash = lockOrder.TxHash
			// } else if lockOrder.Status == lockpaymentorder.StatusPending {
			// 	status = paymentorder.StatusPending
			// 	blockNumber = lockOrder.BlockNumber
			// 	percentSettled = decimal.NewFromInt(0)
			// 	txHash = lockOrder.TxHash
			// }

			// _, err := order.Update().
			// 	SetStatus(status).
			// 	SetBlockNumber(blockNumber).
			// 	SetPercentSettled(percentSettled).
			// 	SetTxHash(txHash).
			// 	Save(ctx)
			// if err != nil {
			// 	logger.Errorf("FixDatabaseMisHap: %v", err)
			// 	continue
			// }

		}
	}

	return nil
}

// HandleReceiveAddressValidity handles receive address validity
func HandleReceiveAddressValidity() error {
	ctx := context.Background()

	// Establish RPC connections
	_, err := setRPCClients(ctx)
	if err != nil {
		return fmt.Errorf("HandleReceiveAddressValidity: %w", err)
	}

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
			receiveaddress.HasPaymentOrder(),
		).
		WithPaymentOrder(func(po *ent.PaymentOrderQuery) {
			po.WithToken(func(tq *ent.TokenQuery) {
				tq.WithNetwork()
			})
			po.WithRecipient()
		}).
		All(ctx)
	if err != nil {
		return fmt.Errorf("HandleReceiveAddressValidity: %w", err)
	}

	var indexerService services.Indexer
	for _, address := range addresses {
		if strings.HasPrefix(address.Edges.PaymentOrder.Edges.Token.Edges.Network.Identifier, "tron") {
			indexerService = services.NewIndexerService(orderService.NewOrderTron())
		} else {
			indexerService = services.NewIndexerService(orderService.NewOrderEVM())
		}

		err := indexerService.HandleReceiveAddressValidity(ctx, rpcClients[address.Edges.PaymentOrder.Edges.Token.Edges.Network.Identifier], address, address.Edges.PaymentOrder)
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

	go ReassignStaleOrderRequest(ctx, orderRequestChan)
}

// fetchExternalRate fetches the external rate for a fiat currency
func fetchExternalRate(currency string) (decimal.Decimal, error) {
	supportedCurrencies := []string{"USD", "NGN", "GHS"}
	isSupported := false
	for _, supported := range supportedCurrencies {
		if strings.ToUpper(currency) == supported {
			isSupported = true
			break
		}
	}
	if !isSupported {
		return decimal.Zero, fmt.Errorf("ComputeMarketRate: currency not support")
	}
	// Fetch stable coin rate from third-party API Quidax (USDT)
	res, err := fastshot.NewClient("https://www.quidax.com").
		Config().SetTimeout(30*time.Second).
		Build().GET(fmt.Sprintf("/api/v1/markets/tickers/usdt%s", strings.ToLower(currency))).
		Retry().Set(3, 5*time.Second).
		Send()
	if err != nil {
		return decimal.Zero, fmt.Errorf("ComputeMarketRate: %w", err)
	}

	data, err := utils.ParseJSONResponse(res.RawResponse)
	if err != nil {
		return decimal.Zero, fmt.Errorf("ComputeMarketRate: %w", err)
	}

	price, err := decimal.NewFromString(data["data"].(map[string]interface{})["ticker"].(map[string]interface{})["buy"].(string))
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
		externalRate, err := fetchExternalRate(currency.Code)
		if err != nil {
			return fmt.Errorf("ComputeMarketRate: %w", err)
		}

		// Fetch rates from token configs with fixed conversion rate
		tokenConfigs, err := storage.Client.ProviderOrderToken.
			Query().
			Where(
				providerordertoken.SymbolIn("USDT", "USDC"),
				providerordertoken.ConversionRateTypeEQ(providerordertoken.ConversionRateTypeFixed),
			).
			Select(providerordertoken.FieldFixedConversionRate).
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
		body, err := fastshot.NewClient(attempt.WebhookURL).
			Config().SetTimeout(30*time.Second).
			Header().Add("X-Paycrest-Signature", attempt.Signature).
			Build().POST("").
			Body().AsJSON(attempt.Payload).
			Send()

		if err != nil || (body.StatusCode() >= 205) {
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
				uid, err := uuid.Parse(attempt.Payload["data"].(map[string]interface{})["senderId"].(string))
				if err != nil {
					return fmt.Errorf("RetryFailedWebhookNotifications.FailedExtraction: %w", err)
				}
				profile, err := storage.Client.SenderProfile.
					Query().
					Where(
						senderprofile.IDEQ(uid),
					).
					WithUser().Only(ctx)
				if err != nil {
					return fmt.Errorf("RetryFailedWebhookNotifications.CouldNotFetchProfile: %w", err)
				}

				_, err = services.SendTemplateEmail(types.SendEmailPayload{
					ToAddress: profile.Edges.User.Email,
					DynamicData: map[string]interface{}{
						"first_name": profile.Edges.User.FirstName,
					},
				}, "d-da75eee4966544ad92dcd060421d4e12")

				if err != nil {
					return fmt.Errorf("RetryFailedWebhookNotifications.SendTemplateEmail: %w", err)
				}
			}

			_, err := attemptUpdate.Save(ctx)
			if err != nil {
				return fmt.Errorf("RetryFailedWebhookNotifications: %w", err)
			}

			continue
		}

		// Webhook notification was successful
		_, err = attempt.Update().
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

	// Compute market rate every 3 minutes
	_, err := scheduler.Cron("*/3 * * * *").Do(ComputeMarketRate)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Refresh provision bucket priority queues every X minutes
	_, err = scheduler.Cron(fmt.Sprintf("*/%d * * * *", orderConf.BucketQueueRebuildInterval)).
		Do(priorityQueue.ProcessBucketQueues)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Retry failed webhook notifications every 59 minutes
	_, err = scheduler.Cron("*/59 * * * *").Do(RetryFailedWebhookNotifications)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Reassign pending order requests every 13 minutes
	_, err = scheduler.Cron("*/13 * * * *").Do(ReassignPendingOrders)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Reassign unvalidated order requests every 2 minutes
	_, err = scheduler.Cron("*/2 * * * *").Do(ReassignUnvalidatedLockOrders)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Reassign unfulfilled order requests every 3 minutes
	_, err = scheduler.Cron("*/3 * * * *").Do(ReassignUnfulfilledLockOrders)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Handle receive address validity every 31 minutes
	_, err = scheduler.Cron("*/31 * * * *").Do(HandleReceiveAddressValidity)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Retry stale user operations every 2 minutes
	_, err = scheduler.Cron("*/2 * * * *").Do(RetryStaleUserOperations)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Index blockchain events every 1 minute
	_, err = scheduler.Cron("*/1 * * * *").Do(IndexBlockchainEvents)
	if err != nil {
		logger.Errorf("StartCronJobs: %v", err)
	}

	// Start scheduler
	scheduler.StartAsync()
}
