package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/lockorderfulfillment"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/providerordertoken"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/provisionbucket"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	cryptoUtils "github.com/paycrest/protocol/utils/crypto"
	"github.com/paycrest/protocol/utils/logger"
	tokenUtils "github.com/paycrest/protocol/utils/token"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type PriorityQueueService struct{}

// NewPriorityQueueService creates a new instance of PriorityQueueService
func NewPriorityQueueService() *PriorityQueueService {
	return &PriorityQueueService{}
}

// ProcessBucketQueues creates a priority queue for each bucket and saves it to redis
func (s *PriorityQueueService) ProcessBucketQueues() error {
	ctx := context.Background()

	buckets, err := s.GetProvisionBuckets(ctx)
	if err != nil {
		return fmt.Errorf("ProcessBucketQueues.GetProvisionBuckets: %w", err)
	}

	for _, bucket := range buckets {
		go s.CreatePriorityQueueForBucket(ctx, bucket)
	}

	return nil
}

// GetProvisionBuckets returns a list of buckets with their providers
func (s *PriorityQueueService) GetProvisionBuckets(ctx context.Context) ([]*ent.ProvisionBucket, error) {
	buckets, err := storage.Client.ProvisionBucket.
		Query().
		Select(provisionbucket.FieldMinAmount, provisionbucket.FieldMaxAmount).
		WithProviderProfiles(func(ppq *ent.ProviderProfileQuery) {
			// ppq.WithProviderRating(func(prq *ent.ProviderRatingQuery) {
			// 	prq.Select(providerrating.FieldTrustScore)
			// })
			ppq.Select(providerprofile.FieldID)

			// Filter only providers that are always available
			ppq.Where(
				providerprofile.IsAvailableEQ(true),
				providerprofile.IsActiveEQ(true),
				providerprofile.VisibilityModeEQ(providerprofile.VisibilityModePublic),
			)
		}).
		WithCurrency().
		All(ctx)
	if err != nil {
		return nil, err
	}

	return buckets, nil
}

// GetProviderRate returns the rate for a provider
func (s *PriorityQueueService) GetProviderRate(ctx context.Context, provider *ent.ProviderProfile) (decimal.Decimal, error) {
	// Fetch the token config for the provider
	tokenConfig, err := storage.Client.ProviderOrderToken.
		Query().
		Where(
			providerordertoken.HasProviderWith(providerprofile.IDEQ(provider.ID)),
		).
		WithProvider(func(pq *ent.ProviderProfileQuery) {
			pq.WithCurrency()
		}).
		Select(
			providerordertoken.FieldConversionRateType,
			providerordertoken.FieldFixedConversionRate,
			providerordertoken.FieldFloatingConversionRate,
		).
		First(ctx)
	if err != nil {
		return decimal.Decimal{}, err
	}

	var rate decimal.Decimal

	if tokenConfig.ConversionRateType == providerordertoken.ConversionRateTypeFixed {
		rate = tokenConfig.FixedConversionRate
	} else {
		// Handle floating rate case
		marketRate := tokenConfig.Edges.Provider.Edges.Currency.MarketRate
		floatingRate := tokenConfig.FloatingConversionRate // in percentage

		// Calculate the floating rate based on the market rate
		deviation := marketRate.Mul(floatingRate.Div(decimal.NewFromInt(100)))
		rate = rate.Add(deviation)
	}

	return rate, nil
}

// CreatePriorityQueueForBucket creates a priority queue for a bucket and saves it to redis
func (s *PriorityQueueService) CreatePriorityQueueForBucket(ctx context.Context, bucket *ent.ProvisionBucket) {
	// Create a slice to store the provider profiles sorted by trust score
	providers := bucket.Edges.ProviderProfiles
	// sort.SliceStable(providers, func(i, j int) bool {
	// 	trustScoreI, _ := providers[i].Edges.ProviderRating.TrustScore.Float64()
	// 	trustScoreJ, _ := providers[j].Edges.ProviderRating.TrustScore.Float64()
	// 	return trustScoreI > trustScoreJ // Sort in descending order
	// })

	// Enqueue provider ID and rate as a single string into the circular queue
	redisKey := fmt.Sprintf("bucket_%s_%s_%s", bucket.Edges.Currency.Code, bucket.MinAmount, bucket.MaxAmount)

	_, err := storage.RedisClient.Del(ctx, redisKey).Result() // delete existing queue
	if err != nil {
		logger.Errorf("failed to delete existing circular queue: %v", err)
	}

	for _, provider := range providers {
		providerID := provider.ID
		rate, _ := s.GetProviderRate(ctx, provider)

		// Check provider's rate against the market rate to ensure it's not too far off
		marketRate := bucket.Edges.Currency.MarketRate

		if marketRate.Cmp(decimal.Zero) != 0 {
			if rate.LessThan(marketRate.Mul(decimal.NewFromFloat(1).Sub(OrderConf.PercentDeviationFromMarketRate))) ||
				rate.GreaterThan(marketRate.Mul(decimal.NewFromFloat(1).Add(OrderConf.PercentDeviationFromMarketRate))) {
				// Skip this provider if the rate is too far off
				// TODO: add a logic to notify the provider(s) to update his rate since it's stale. could be a cron job
				continue
			}
		}

		// Serialize the provider ID and rate into a single string
		data := fmt.Sprintf("%s:%s:%s", providerID, rate, strconv.FormatBool(provider.IsPartner))

		// Enqueue the serialized data into the circular queue
		err := storage.RedisClient.RPush(ctx, redisKey, data).Err()
		if err != nil {
			logger.Errorf("failed to enqueue provider data to circular queue: %v", err)
		}
	}
}

// AssignLockPaymentOrders assigns lock payment orders to providers
func (s *PriorityQueueService) AssignLockPaymentOrder(ctx context.Context, order types.LockPaymentOrderFields) error {
	s.ReassignUnfulfilledLockOrders(ctx)

	excludeList, err := storage.RedisClient.LRange(ctx, fmt.Sprintf("order_exclude_list_%s", order.ID), 0, -1).Result()
	if err != nil {
		logger.Errorf("failed to get exclude list for order %d: %v", order.ID, err)
		return err
	}

	// Sends order directly to the specified provider in order. Incase of failure, proceed to queue
	if order.ProviderID != "" && !utils.ContainsString(excludeList, order.ProviderID) {
		err := s.sendOrderRequest(ctx, order)
		if err == nil {
			return nil
		}
		logger.Errorf("failed to send order request to specific provider %s: %v. sending order to queue",
			order.ProviderID, err)
	}

	// Get the first provider from the circular queue
	redisKey := fmt.Sprintf("bucket_%s_%s_%s", order.ProvisionBucket.Edges.Currency.Code, order.ProvisionBucket.MinAmount, order.ProvisionBucket.MaxAmount)

	partnerProviders := []string{}

	for index := 0; ; index++ {
		providerData, err := storage.RedisClient.LIndex(ctx, redisKey, int64(index)).Result()
		if err != nil {
			logger.Errorf("failed to access index %d from circular queue: %v", index, err)
			break
		}

		if providerData == "" {
			// Reached the end of the queue
			logger.Errorf("rate didn't match a provider, finding a partner provider")

			if len(partnerProviders) == 0 {
				logger.Errorf("no partner providers found")
				return nil
			}

			// Pick a random partner provider
			randomIndex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(partnerProviders))
			providerData = partnerProviders[randomIndex]
		}

		// Extract the rate from the data (assuming it's in the format "providerID:rate:is_partner")
		parts := strings.Split(providerData, ":")
		if len(parts) != 3 {
			logger.Errorf("invalid data format at index %d: %s", index, providerData)
			continue // Skip this entry due to invalid format
		}

		// Check if the provider is a partner
		if parts[2] == "true" {
			partnerProviders = append(partnerProviders, providerData)
		}

		order.ProviderID = parts[0]

		// Skip entry if provider is excluded
		if utils.ContainsString(excludeList, order.ProviderID) {
			continue
		}

		rate, err := decimal.NewFromString(parts[1])
		if err != nil {
			logger.Errorf("failed to parse rate at index %d: %v", index, err)
			continue // Skip this entry due to parsing error
		}

		if rate.Equal(order.Rate) {
			// Found a match for the rate
			if index == 0 {
				// Match found at index 0, perform LPOP to dequeue
				data, err := storage.RedisClient.LPop(ctx, redisKey).Result()
				if err != nil {
					logger.Errorf("failed to dequeue from circular queue: %v", err)
					return err
				}

				// Enqueue data to the end of the queue
				err = storage.RedisClient.RPush(ctx, redisKey, data).Err()
				if err != nil {
					logger.Errorf("failed to enqueue to circular queue: %v", err)
					return err
				}
			}

			// Assign the order to the provider and save it to Redis
			err = s.sendOrderRequest(ctx, order)
			if err != nil {
				logger.Errorf("failed to send order request to specific provider %s: %v", order.ProviderID, err)

				// Push provider ID to order exclude list
				orderKey := fmt.Sprintf("order_exclude_list_%s", order.ID)
				_, err = storage.RedisClient.RPush(ctx, orderKey, order.ProviderID).Result()
				if err != nil {
					logger.Errorf("error pushing provider %s to order %d exclude_list on Redis: %v", order.ProviderID, order.ID, err)
				}

				// Reassign the lock payment order to another provider
				return s.AssignLockPaymentOrder(ctx, order)
			}

			break
		}
	}

	return nil
}

// sendOrderRequest sends an order request to a provider
func (s *PriorityQueueService) sendOrderRequest(ctx context.Context, order types.LockPaymentOrderFields) error {

	// Assign the order to the provider and save it to Redis
	orderKey := fmt.Sprintf("order_request_%s", order.ID)

	approxAmount := order.Amount.Mul(order.Rate).Floor()
	approxAmount = approxAmount.Round(2)

	orderRequestData := map[string]interface{}{
		"amount":      approxAmount.String(),
		"institution": order.Institution,
		"providerId":  order.ProviderID,
	}

	if err := storage.RedisClient.HSet(ctx, orderKey, orderRequestData).Err(); err != nil {
		logger.Errorf("failed to map order to a provider in Redis: %v", err)
		return err
	}

	// Set a TTL for the order request
	err := storage.RedisClient.ExpireAt(ctx, orderKey, time.Now().Add(OrderConf.OrderRequestValidity)).Err()
	if err != nil {
		logger.Errorf("failed to set TTL for order request: %v", err)
		return err
	}

	// Notify the provider
	orderRequestData["orderId"] = order.ID
	if err := s.notifyProvider(ctx, orderRequestData); err != nil {
		logger.Errorf("failed to notify provider %s: %v", order.ProviderID, err)
		return err
	}

	return nil
}

// notifyProvider sends an order request notification to a provider
// TODO: ideally notifications should be moved to a notification service
func (s *PriorityQueueService) notifyProvider(ctx context.Context, orderRequestData map[string]interface{}) error {
	// TODO: can we add mode and host identifier to redis during priority queue creation?
	providerID := orderRequestData["providerId"].(string)
	delete(orderRequestData, "providerId")

	provider, err := storage.Client.ProviderProfile.
		Query().
		Where(
			providerprofile.IDEQ(providerID),
		).
		WithAPIKey().
		Select(providerprofile.FieldProvisionMode, providerprofile.FieldHostIdentifier).
		Only(ctx)
	if err != nil {
		return err
	}

	if provider.ProvisionMode == providerprofile.ProvisionModeAuto {
		// Compute HMAC
		decodedSecret, err := base64.StdEncoding.DecodeString(provider.Edges.APIKey.Secret)
		if err != nil {
			return err
		}
		decryptedSecret, err := cryptoUtils.DecryptPlain(decodedSecret)
		if err != nil {
			return err
		}
		signature := tokenUtils.GenerateHMACSignature(orderRequestData, string(decryptedSecret))

		// Send POST request to the provider's node
		_, err = utils.MakeJSONRequest(
			ctx,
			"POST",
			fmt.Sprintf("%s/new_order", provider.HostIdentifier),
			orderRequestData,
			map[string]string{
				"X-Paycrest-Signature": signature,
			},
		)
		if err != nil {
			return err
		}
	} else {
		notificationConf := config.NotificationConfig()
		// Send POST request to the provider's mobile device
		requestBody := map[string]interface{}{
			"app_id": notificationConf.OneSignalAppID,
			"include_aliases": map[string]interface{}{
				"external_id": []string{providerID},
			},
			"target_channel": "push",
			"headings": map[string]interface{}{
				"en": "Incoming Order 💸",
			},
			"contents": map[string]interface{}{
				// TODO: format currency string with commas
				"en": "You have a payment order of ₦" + orderRequestData["amount"].(decimal.Decimal).String(),
			},
			"data": orderRequestData,
			"buttons": []map[string]interface{}{
				{"id": "accept-button", "text": "Accept"},
				{"id": "decline-button", "text": "Decline"},
			},
			"name": "Order Requests",
		}

		headers := map[string]string{
			"accept":        "application/json",
			"Authorization": fmt.Sprintf("Basic %s", notificationConf.OneSignalRESTAPIKey),
		}
		_, err := utils.MakeJSONRequest(
			ctx,
			"POST",
			"https://onesignal.com/api/v1/notifications",
			requestBody,
			headers,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReassignStaleOrderRequest reassigns expired order requests to providers
func (s *PriorityQueueService) ReassignStaleOrderRequest(ctx context.Context, orderRequestChan <-chan *redis.Message) {
	for msg := range orderRequestChan {
		key := strings.Split(msg.Payload, "_")
		orderID := key[len(key)-1]

		orderUUID, err := uuid.Parse(orderID)
		if err != nil {
			logger.Errorf("ReassignStaleOrderRequest: %v", err)
			return
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
			return
		}

		orderFields := types.LockPaymentOrderFields{
			ID:                order.ID,
			OrderID:           order.OrderID,
			Amount:            order.Amount,
			Rate:              order.Rate,
			Label:             order.Label,
			BlockNumber:       order.BlockNumber,
			Institution:       order.Institution,
			AccountIdentifier: order.AccountIdentifier,
			AccountName:       order.AccountName,
			Memo:              order.Memo,
			ProvisionBucket:   order.Edges.ProvisionBucket,
		}

		// Assign the order to a provider
		err = s.AssignLockPaymentOrder(ctx, orderFields)
		if err != nil {
			logger.Errorf("ReassignStaleOrderRequest: %v", err)
			return
		}
	}
}

// ReassignUnfulfilledLockOrders reassigns lockOrder unfulfilled within a time frame.
func (s *PriorityQueueService) ReassignUnfulfilledLockOrders(ctx context.Context) {
	// Query unfulfilled lock orders.
	lockOrders, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.Not(lockpaymentorder.HasFulfillment()),
			lockpaymentorder.Or(
				lockpaymentorder.And(
					lockpaymentorder.StatusEQ(lockpaymentorder.StatusProcessing),
					lockpaymentorder.UpdatedAtLTE(time.Now().Add(-OrderConf.OrderFulfillmentValidity*time.Minute)),
				),
				lockpaymentorder.StatusEQ(lockpaymentorder.StatusCancelled),
			),
		).
		WithToken().
		WithProvider().
		WithProvisionBucket().
		All(ctx)
	if err != nil {
		logger.Errorf("failed to fetch unfulfilled lock order: %v", err)
		return
	}

	// Unassign unfulfilled lock orders.
	_, err = storage.Client.LockPaymentOrder.
		Update().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusProcessing),
			lockpaymentorder.UpdatedAtLTE(time.Now().Add(-OrderConf.OrderFulfillmentValidity*time.Minute)),
			lockpaymentorder.Not(lockpaymentorder.HasFulfillment()),
		).
		SetStatus(lockpaymentorder.StatusPending).
		Save(ctx)
	if err != nil {
		logger.Errorf("failed to unassign unfulfilled lock order: %v", err)
		return
	}

	for _, order := range lockOrders {
		lockPaymentOrder := types.LockPaymentOrderFields{
			ID:                order.ID,
			Token:             order.Edges.Token,
			OrderID:           order.OrderID,
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

		err := s.AssignLockPaymentOrder(ctx, lockPaymentOrder)
		if err != nil {
			logger.Errorf("failed to reassign unfulfilled lock order with id: %s => %v", order.OrderID, err)
		}
	}
}

// ReassignUnvalidatedLockOrders reassigns unvalidated lock orders to providers
func (s *PriorityQueueService) ReassignUnvalidatedLockOrders() {
	ctx := context.Background()

	// Query unvalidated lock orders.
	lockOrders, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusFulfilled),
			lockpaymentorder.HasFulfillmentWith(
				lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusFailed),
			),
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
		logger.Errorf("ReassignUnvalidatedLockOrders.db: %v", err)
		return
	}

	for _, order := range lockOrders {
		lockPaymentOrder := types.LockPaymentOrderFields{
			ID:                order.ID,
			Token:             order.Edges.Token,
			OrderID:           order.OrderID,
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

		err := s.AssignLockPaymentOrder(ctx, lockPaymentOrder)
		if err != nil {
			logger.Errorf("failed to reassign unvalidated order request: %v", err)
		}
	}
}

// ReassignPendingOrders reassigns declined order requests to providers
func (s *PriorityQueueService) ReassignPendingOrders() {
	ctx := context.Background()

	// Query pending lock orders
	lockOrders, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusPending),
			lockpaymentorder.Not(lockpaymentorder.HasFulfillment()),
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
			return
		}

		if exists == 0 {
			// Order request doesn't exist in Redis, reassign the order
			lockPaymentOrder := types.LockPaymentOrderFields{
				ID:                order.ID,
				Token:             order.Edges.Token,
				OrderID:           order.OrderID,
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

			err := s.AssignLockPaymentOrder(ctx, lockPaymentOrder)
			if err != nil {
				logger.Errorf("failed to reassign declined order request: %v", err)
			}
		}
	}
}
