package provider

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/paycrest/aggregator/config"
	"github.com/paycrest/aggregator/ent"
	"github.com/paycrest/aggregator/ent/fiatcurrency"
	"github.com/paycrest/aggregator/ent/lockorderfulfillment"
	"github.com/paycrest/aggregator/ent/lockpaymentorder"
	"github.com/paycrest/aggregator/ent/providerprofile"
	"github.com/paycrest/aggregator/ent/token"
	"github.com/paycrest/aggregator/ent/transactionlog"
	orderService "github.com/paycrest/aggregator/services/order"
	"github.com/paycrest/aggregator/storage"
	"github.com/paycrest/aggregator/types"
	u "github.com/paycrest/aggregator/utils"
	"github.com/paycrest/aggregator/utils/logger"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

var orderConf = config.OrderConfig()

// ProviderController is a controller type for provider endpoints
type ProviderController struct{}

// NewProviderController creates a new instance of ProviderController with injected services
func NewProviderController() *ProviderController {
	return &ProviderController{}
}

// GetLockPaymentOrders controller fetches all assigned orders
func (ctrl *ProviderController) GetLockPaymentOrders(ctx *gin.Context) {
	// get page and pageSize query params
	page, offset, pageSize := u.Paginate(ctx)

	// Set ordering
	ordering := ctx.Query("ordering")
	order := ent.Desc(lockpaymentorder.FieldCreatedAt)
	if ordering == "asc" {
		order = ent.Asc(lockpaymentorder.FieldCreatedAt)
	}

	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	lockPaymentOrderQuery := storage.Client.LockPaymentOrder.Query()

	// Filter by status
	statusMap := map[string]lockpaymentorder.Status{
		"pending":    lockpaymentorder.StatusPending,
		"validated":  lockpaymentorder.StatusValidated,
		"fulfilled":  lockpaymentorder.StatusFulfilled,
		"cancelled":  lockpaymentorder.StatusCancelled,
		"processing": lockpaymentorder.StatusProcessing,
		"settled":    lockpaymentorder.StatusSettled,
	}

	statusQueryParam := ctx.Query("status")

	if status, ok := statusMap[statusQueryParam]; ok {
		lockPaymentOrderQuery = lockPaymentOrderQuery.Where(
			lockpaymentorder.HasProviderWith(providerprofile.IDEQ(provider.ID)),
			lockpaymentorder.StatusEQ(status),
		)
	} else {
		lockPaymentOrderQuery = lockPaymentOrderQuery.Where(
			lockpaymentorder.HasProviderWith(providerprofile.IDEQ(provider.ID)),
		)
	}

	count, err := lockPaymentOrderQuery.Count(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch orders", nil)
		return
	}

	// Fetch all orders assigned to the provider
	lockPaymentOrders, err := lockPaymentOrderQuery.
		Limit(pageSize).
		Offset(offset).
		Order(order).
		WithProvider().
		WithToken(
			func(query *ent.TokenQuery) {
				query.WithNetwork()
			},
		).
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch orders", nil)
		return
	}

	var orders []types.LockPaymentOrderResponse
	for _, order := range lockPaymentOrders {
		orders = append(orders, types.LockPaymentOrderResponse{
			ID:                  order.ID,
			Token:               order.Edges.Token.Symbol,
			GatewayID:           order.GatewayID,
			Amount:              order.Amount,
			Rate:                order.Rate,
			Institution:         order.Institution,
			AccountIdentifier:   order.AccountIdentifier,
			AccountName:         order.AccountName,
			TxHash:              order.TxHash,
			Status:              order.Status,
			Memo:                order.Memo,
			Network:             order.Edges.Token.Edges.Network.Identifier,
			CancellationReasons: order.CancellationReasons,
			UpdatedAt:           order.UpdatedAt,
			CreatedAt:           order.CreatedAt,
		})
	}

	// return paginated orders
	u.APIResponse(ctx, http.StatusOK, "success", "Orders successfully retrieved", types.ProviderLockOrderList{
		Page:         page,
		PageSize:     pageSize,
		TotalRecords: count,
		Orders:       orders,
	})
}

// AcceptOrder controller accepts an order
func (ctrl *ProviderController) AcceptOrder(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Parse the Order ID string into a UUID
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logger.Errorf("error parsing order ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid Order ID", nil)
		return
	}

	// Get Order request from Redis
	result, err := storage.RedisClient.HGetAll(ctx, fmt.Sprintf("order_request_%s", orderID)).Result()
	if err != nil {
		logger.Errorf("error getting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to accept order request", nil)
		return
	}

	if result["providerId"] != provider.ID || len(result) == 0 {
		logger.Errorf("order request not found in Redis: %v", orderID)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Order request not found or is expired", nil)
		return
	}

	// Delete order request from Redis
	_, err = storage.RedisClient.Del(ctx, fmt.Sprintf("order_request_%s", orderID)).Result()
	if err != nil {
		logger.Errorf("error deleting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to accept order request", nil)
		return
	}

	tx, err := storage.Client.Tx(ctx)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		return
	}

	// Log transaction status
	var transactionLog *ent.TransactionLog
	_, err = tx.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.IDEQ(orderID),
			lockpaymentorder.HasTransactionsWith(
				transactionlog.StatusEQ(transactionlog.StatusOrderProcessing),
			),
		).
		Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		} else {
			transactionLog, err = tx.TransactionLog.
				Create().
				SetStatus(transactionlog.StatusOrderProcessing).
				SetMetadata(
					map[string]interface{}{
						"ProviderId": provider.ID,
					}).
				Save(ctx)
			if err != nil {
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
				return
			}
		}
	}

	// Update lock order status to processing
	orderBuilder := tx.LockPaymentOrder.
		UpdateOneID(orderID).
		SetStatus(lockpaymentorder.StatusProcessing).
		SetProviderID(provider.ID)

	if transactionLog != nil {
		orderBuilder = orderBuilder.AddTransactions(transactionLog)
	}

	order, err := orderBuilder.Save(ctx)
	if err != nil {
		logger.Errorf("%s - error.AcceptOrder: %v", orderID, err)
		if ent.IsNotFound(err) {
			u.APIResponse(ctx, http.StatusNotFound, "error", "Order not found", nil)
		} else {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		}
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		return
	}

	u.APIResponse(ctx, http.StatusCreated, "success", "Order request accepted successfully", &types.AcceptOrderResponse{
		ID:                orderID,
		Amount:            order.Amount.Mul(order.Rate).RoundBank(0),
		Institution:       order.Institution,
		AccountIdentifier: order.AccountIdentifier,
		AccountName:       order.AccountName,
		Memo:              order.Memo,
	})
}

// DeclineOrder controller declines an order
func (ctrl *ProviderController) DeclineOrder(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Parse the Order ID string into a UUID
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logger.Errorf("error parsing order ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid Order ID", nil)
		return
	}

	// Get Order request from Redis
	result, err := storage.RedisClient.HGetAll(ctx, fmt.Sprintf("order_request_%s", orderID)).Result()
	if err != nil {
		logger.Errorf("error getting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decline order request", nil)
		return
	}

	if result["providerId"] != provider.ID || len(result) == 0 {
		logger.Errorf("order request not found in Redis: %v", orderID)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Order request not found or is expired", nil)
		return
	}

	// Delete order request from Redis
	_, err = storage.RedisClient.Del(ctx, fmt.Sprintf("order_request_%s", orderID)).Result()
	if err != nil {
		logger.Errorf("error deleting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decline order request", nil)
		return
	}

	// Push provider ID to order exclude list
	orderKey := fmt.Sprintf("order_exclude_list_%s", orderID)
	_, err = storage.RedisClient.RPush(ctx, orderKey, provider.ID).Result()
	if err != nil {
		logger.Errorf("error pushing provider %s to order %s exclude_list on Redis: %v", provider.ID, orderID, err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decline order request", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Order request declined successfully", nil)
}

// FulfillOrder controller fulfills an order
func (ctrl *ProviderController) FulfillOrder(ctx *gin.Context) {
	var payload types.FulfillLockOrderPayload

	// Parse the order payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Get provider profile from the context
	_, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}

	// Parse the Order ID string into a UUID
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logger.Errorf("error parsing order ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid Order ID", nil)
		return
	}

	updateLockOrder := storage.Client.LockPaymentOrder.
		Update().
		Where(
			lockpaymentorder.IDEQ(orderID),
			lockpaymentorder.Or(
				lockpaymentorder.StatusEQ(lockpaymentorder.StatusProcessing),
				lockpaymentorder.StatusEQ(lockpaymentorder.StatusFulfilled),
			),
		)

	// Query or create lock order fulfillment
	fulfillment, err := storage.Client.LockOrderFulfillment.
		Query().
		Where(lockorderfulfillment.TxIDEQ(payload.TxID)).
		WithOrder(func(poq *ent.LockPaymentOrderQuery) {
			poq.WithToken(func(tq *ent.TokenQuery) {
				tq.WithNetwork()
			})
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			_, err = storage.Client.LockOrderFulfillment.
				Create().
				SetOrderID(orderID).
				SetTxID(payload.TxID).
				SetPsp(payload.PSP).
				Save(ctx)
			if err != nil {
				logger.Errorf("error: %v", err)
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
				return
			}

			fulfillment, err = storage.Client.LockOrderFulfillment.
				Query().
				Where(lockorderfulfillment.TxIDEQ(payload.TxID)).
				WithOrder(func(poq *ent.LockPaymentOrderQuery) {
					poq.WithToken(func(tq *ent.TokenQuery) {
						tq.WithNetwork()
					})
				}).
				Only(ctx)
			if err != nil {
				logger.Errorf("error: %v", err)
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
				return
			}
		} else {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}
	}

	if payload.ValidationStatus == lockorderfulfillment.ValidationStatusSuccess {
		if fulfillment.Edges.Order.Status != lockpaymentorder.StatusFulfilled {
			u.APIResponse(ctx, http.StatusOK, "success", "Order already validated", nil)
			return
		}

		_, err := fulfillment.Update().
			SetValidationStatus(lockorderfulfillment.ValidationStatusSuccess).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}

		transactionLog, err := storage.Client.TransactionLog.Create().
			SetStatus(transactionlog.StatusOrderValidated).
			SetNetwork(fulfillment.Edges.Order.Edges.Token.Edges.Network.Identifier).
			SetMetadata(map[string]interface{}{
				"TransactionID": payload.TxID,
				"PSP":           payload.PSP,
			}).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}

		_, err = updateLockOrder.
			SetStatus(lockpaymentorder.StatusValidated).
			AddTransactions(transactionLog).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}

		// Settle order or fail silently
		go func() {
			var err error
			if strings.HasPrefix(fulfillment.Edges.Order.Edges.Token.Edges.Network.Identifier, "tron") {
				err = orderService.NewOrderTron().SettleOrder(ctx, nil, orderID)
			} else {
				err = orderService.NewOrderEVM().SettleOrder(ctx, nil, orderID)
			}
			if err != nil {
				logger.Errorf("FulfillOrder.SettleOrder: %v", err)
			}
		}()

	} else if payload.ValidationStatus == lockorderfulfillment.ValidationStatusFailed {
		_, err = fulfillment.Update().
			SetValidationStatus(lockorderfulfillment.ValidationStatusFailed).
			SetValidationError(payload.ValidationError).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}

		_, err = updateLockOrder.
			SetStatus(lockpaymentorder.StatusFulfilled).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}

	} else {
		transactionLog, err := storage.Client.TransactionLog.Create().
			SetStatus(transactionlog.StatusOrderFulfilled).
			SetNetwork(fulfillment.Edges.Order.Edges.Token.Edges.Network.Identifier).
			SetMetadata(map[string]interface{}{
				"TransactionID": payload.TxID,
				"PSP":           payload.PSP,
			}).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}

		_, err = updateLockOrder.
			SetStatus(lockpaymentorder.StatusFulfilled).
			AddTransactions(transactionLog).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Order fulfilled successfully", nil)
}

// CancelOrder controller cancels an order
func (ctrl *ProviderController) CancelOrder(ctx *gin.Context) {
	var payload types.CancelLockOrderPayload

	// Parse the order payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Parse the Order ID string into a UUID
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logger.Errorf("error parsing order ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid Order ID", nil)
		return
	}

	// Fetch lock payment order from db
	order, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.IDEQ(orderID),
			lockpaymentorder.HasProviderWith(providerprofile.IDEQ(provider.ID)),
		).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		WithProvider().
		WithProvisionBucket(func(pbq *ent.ProvisionBucketQuery) {
			pbq.WithCurrency()
		}).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Could not find payment order", nil)
		return
	}

	// Get new cancellation count based on cancel reason
	orderUpdate := storage.Client.LockPaymentOrder.UpdateOneID(orderID)
	cancellationCount := order.CancellationCount
	if payload.Reason == "Invalid recipient bank details" || provider.VisibilityMode == providerprofile.VisibilityModePrivate {
		cancellationCount += orderConf.RefundCancellationCount // Allows us refund immediately for invalid recipient
		orderUpdate.AppendCancellationReasons([]string{payload.Reason})
	} else if payload.Reason != "Insufficient funds" {
		cancellationCount += 1
		orderUpdate.AppendCancellationReasons([]string{payload.Reason})
	} else if payload.Reason == "Insufficient funds" {
		// Search for the specific provider in the queue using a Redis list
		redisKey := fmt.Sprintf("bucket_%s_%s_%s", order.Edges.ProvisionBucket.Edges.Currency.Code, order.Edges.ProvisionBucket.MinAmount, order.Edges.ProvisionBucket.MaxAmount)

		// Check if the provider ID exists in the list
		for index := -1; ; index-- {
			providerData, err := storage.RedisClient.LIndex(ctx, redisKey, int64(index)).Result()
			if err != nil {
				break
			}

			// Extract the id from the data (assuming format "providerID:token:rate:minAmount:maxAmount")
			parts := strings.Split(providerData, ":")
			if len(parts) != 5 {
				logger.Errorf("invalid provider data format: %s", providerData)
				continue // Skip this entry due to invalid format
			}

			if parts[0] == provider.ID {
				// Remove the provider from the list
				placeholder := "DELETED_PROVIDER" // Define a placeholder value
				_, err := storage.RedisClient.LSet(ctx, redisKey, int64(index), placeholder).Result()
				if err != nil {
					logger.Errorf("failed to set placeholder at index %d: %v", index, err)
				}

				// Remove all occurences of the placeholder from the list
				_, err = storage.RedisClient.LRem(ctx, redisKey, 0, placeholder).Result()
				if err != nil {
					logger.Errorf("failed to remove placeholder from circular queue: %v", err)
				}

				break
			}
		}

		// // Update provider availability to off
		// _, err = storage.Client.ProviderProfile.
		// 	UpdateOneID(provider.ID).
		// 	SetIsAvailable(false).
		// 	Save(ctx)
		// if err != nil {
		// 	logger.Errorf("failed to update provider availability: %v", err)
		// }
	}

	// Update lock order status to cancelled
	_, err = orderUpdate.
		SetStatus(lockpaymentorder.StatusCancelled).
		SetCancellationCount(cancellationCount).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to cancel order", nil)
		return
	}

	order.Status = lockpaymentorder.StatusCancelled
	order.CancellationCount = cancellationCount

	// Check if order cancellation count is equal or greater than RefundCancellationCount in config,
	// and the order has not been refunded, then trigger refund
	if order.CancellationCount >= orderConf.RefundCancellationCount && order.Status == lockpaymentorder.StatusCancelled {
		go func() {
			var err error
			if strings.HasPrefix(order.Edges.Token.Edges.Network.Identifier, "tron") {
				err = orderService.NewOrderTron().RefundOrder(ctx, nil, order.Edges.Token.Edges.Network, order.GatewayID)
			} else {
				err = orderService.NewOrderEVM().RefundOrder(ctx, nil, order.Edges.Token.Edges.Network, order.GatewayID)
			}
			if err != nil {
				logger.Errorf("CancelOrder.RefundOrder(%v): %v", orderID, err)
			}
		}()
	}

	// Push provider ID to order exclude list
	orderKey := fmt.Sprintf("order_exclude_list_%s", orderID)
	_, err = storage.RedisClient.RPush(ctx, orderKey, provider.ID).Result()
	if err != nil {
		logger.Errorf("error pushing provider %s to order %s exclude_list on Redis: %v", provider.ID, orderID, err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decline order request", nil)
		return
	}

	// TODO: Reassign order to another provider in background

	u.APIResponse(ctx, http.StatusOK, "success", "Order cancelled successfully", nil)
}

// GetMarketRate controller fetches the median rate of the cryptocurrency token against the fiat currency
func (ctrl *ProviderController) GetMarketRate(ctx *gin.Context) {
	// Parse path parameters
	tokenExists, err := storage.Client.Token.
		Query().
		Where(
			token.SymbolEQ(strings.ToUpper(ctx.Param("token"))),
			token.IsEnabledEQ(true),
		).
		Exist(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to get market rate", nil)
		return
	}

	if !tokenExists {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Token is not supported", nil)
		return
	}
	// TODO: use token to get the token rate for that currency based on the USD/Token Ratio USD/USDC can be 1.005 and USD/USD can be 0.9995

	currency, err := storage.Client.FiatCurrency.
		Query().
		Where(
			fiatcurrency.IsEnabledEQ(true),
			fiatcurrency.CodeEQ(strings.ToUpper(ctx.Param("fiat"))),
		).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Fiat currency is not supported", nil)
		return
	}

	deviation := currency.MarketRate.Mul(orderConf.PercentDeviationFromMarketRate.Div(decimal.NewFromInt(100)))

	u.APIResponse(ctx, http.StatusOK, "success", "Rate fetched successfully", &types.MarketRateResponse{
		MarketRate:  currency.MarketRate,
		MinimumRate: currency.MarketRate.Sub(deviation),
		MaximumRate: currency.MarketRate.Add(deviation),
	})
}

// Stats controller fetches provider stats
func (ctrl *ProviderController) Stats(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Fetch provider stats
	query := storage.Client.LockPaymentOrder.
		Query().
		Where(lockpaymentorder.HasProviderWith(providerprofile.IDEQ(provider.ID)), lockpaymentorder.StatusEQ(lockpaymentorder.StatusSettled))

	var v []struct {
		Sum decimal.Decimal
	}

	err := query.
		Aggregate(
			ent.Sum(lockpaymentorder.FieldAmount),
		).
		Scan(ctx, &v)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch provider stats", nil)
		return
	}

	settledOrders, err := query.
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch provider stats", nil)
		return
	}

	var totalFiatVolume decimal.Decimal
	for _, order := range settledOrders {
		totalFiatVolume = totalFiatVolume.Add(order.Amount.Mul(order.Rate).RoundBank(0))
	}

	count, err := storage.Client.LockPaymentOrder.
		Query().
		Where(lockpaymentorder.HasProviderWith(providerprofile.IDEQ(provider.ID))).
		Count(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch provider stats", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Provider stats fetched successfully", &types.ProviderStatsResponse{
		TotalOrders:       count,
		TotalFiatVolume:   totalFiatVolume,
		TotalCryptoVolume: v[0].Sum,
	})
}

// NodeInfo controller fetches the provision node info
func (ctrl *ProviderController) NodeInfo(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}

	provider, err := storage.Client.ProviderProfile.
		Query().
		Where(providerprofile.IDEQ(providerCtx.(*ent.ProviderProfile).ID)).
		WithAPIKey().
		WithCurrency().
		Only(ctx)
	if err != nil {
		logger.Errorf("failed to fetch provider: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch node info", nil)
		return
	}

	res, err := fastshot.NewClient(provider.HostIdentifier).
		Config().SetTimeout(30 * time.Second).
		Build().GET("/health").
		Send()
	if err != nil {
		logger.Errorf("failed to fetch node info: %v", err)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to fetch node info", nil)
		return
	}

	data, err := u.ParseJSONResponse(res.RawResponse)
	if err != nil {
		logger.Errorf("failed to parse node info: %v", err)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to fetch node info", nil)
		return
	}

	currency := data["data"].(map[string]interface{})["currency"].(string)
	if currency != provider.Edges.Currency.Code {
		logger.Errorf("failed to match currency: %v", err)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to fetch node info", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Node info fetched successfully", data)
}

// GetLockPaymentOrderByID controller fetches a payment order by ID
func (ctrl *ProviderController) GetLockPaymentOrderByID(ctx *gin.Context) {
	// Get order ID from the URL
	orderID := ctx.Param("id")

	// Convert order ID to UUID
	id, err := uuid.Parse(orderID)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Invalid order ID", nil)
		return
	}

	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")

	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Fetch payment order from the database
	lockPaymentOrder, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.IDEQ(id),
			lockpaymentorder.HasProviderWith(providerprofile.IDEQ(provider.ID)),
		).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		WithTransactions().
		Only(ctx)

	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error",
			"Payment order not found", nil)
		return
	}
	var transactions []types.TransactionLog
	for _, transaction := range lockPaymentOrder.Edges.Transactions {
		transactions = append(transactions, types.TransactionLog{
			ID:        transaction.ID,
			GatewayId: transaction.GatewayID,
			Status:    transaction.Status,
			TxHash:    transaction.TxHash,
			CreatedAt: transaction.CreatedAt,
		})

	}

	u.APIResponse(ctx, http.StatusOK, "success", "The order has been successfully retrieved", &types.LockPaymentOrderResponse{
		ID:                  lockPaymentOrder.ID,
		Token:               lockPaymentOrder.Edges.Token.Symbol,
		GatewayID:           lockPaymentOrder.GatewayID,
		Amount:              lockPaymentOrder.Amount,
		Rate:                lockPaymentOrder.Rate,
		Institution:         lockPaymentOrder.Institution,
		AccountIdentifier:   lockPaymentOrder.AccountIdentifier,
		AccountName:         lockPaymentOrder.AccountName,
		TxHash:              lockPaymentOrder.TxHash,
		Status:              lockPaymentOrder.Status,
		Memo:                lockPaymentOrder.Memo,
		Network:             lockPaymentOrder.Edges.Token.Edges.Network.Identifier,
		UpdatedAt:           lockPaymentOrder.UpdatedAt,
		CreatedAt:           lockPaymentOrder.CreatedAt,
		Transactions:        transactions,
		CancellationReasons: lockPaymentOrder.CancellationReasons,
	})
}
