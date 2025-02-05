package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/paycrest/aggregator/ent"
	"github.com/paycrest/aggregator/ent/fiatcurrency"
	"github.com/paycrest/aggregator/ent/paymentorder"
	"github.com/paycrest/aggregator/ent/providerordertoken"
	"github.com/paycrest/aggregator/ent/providerprofile"
	"github.com/paycrest/aggregator/ent/provisionbucket"
	"github.com/paycrest/aggregator/ent/token"
	"github.com/paycrest/aggregator/storage"
	"github.com/paycrest/aggregator/types"
	"github.com/paycrest/aggregator/utils"
	cryptoUtils "github.com/paycrest/aggregator/utils/crypto"
	"github.com/paycrest/aggregator/utils/logger"
	tokenUtils "github.com/paycrest/aggregator/utils/token"
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
				providerprofile.IsAvailable(true),
				providerprofile.IsActive(true),
				providerprofile.IsKybVerified(true),
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
func (s *PriorityQueueService) GetProviderRate(ctx context.Context, provider *ent.ProviderProfile, tokenSymbol string, currency string) (decimal.Decimal, error) {
	// Fetch the token config for the provider
	tokenConfig, err := storage.Client.ProviderOrderToken.
		Query().
		Where(
			providerordertoken.HasProviderWith(providerprofile.IDEQ(provider.ID)),
			providerordertoken.HasTokenWith(token.SymbolEQ(tokenSymbol)),
			providerordertoken.HasCurrencyWith(fiatcurrency.CodeEQ(currency)),
		).
		WithProvider().
		WithCurrency().
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
		marketRate := tokenConfig.Edges.Currency.MarketRate
		floatingRate := tokenConfig.FloatingConversionRate // in percentage

		// Calculate the floating rate based on the market rate
		rate = marketRate.Add(floatingRate).RoundBank(2)
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

	// TODO: add also the checks for all the currencies that a provider has

	for _, provider := range providers {
		orderTokens, err := storage.Client.ProviderOrderToken.
			Query().
			Where(
				providerordertoken.HasProviderWith(providerprofile.IDEQ(provider.ID)),
				providerordertoken.HasCurrencyWith(fiatcurrency.CodeEQ(bucket.Edges.Currency.Code)),
			).
			WithToken().
			All(ctx)
		if err != nil {
			fmt.Println("error", err)
			logger.Errorf("failed to get tokens for provider %s: %v", provider.ID, err)
			continue
		}

		fmt.Println("orderTokens", orderTokens)

		for _, orderToken := range orderTokens {
			fmt.Println("orderToken", orderToken)
			providerID := provider.ID
			rate, err := s.GetProviderRate(ctx, provider, orderToken.Edges.Token.Symbol, bucket.Edges.Currency.Code)
			if err != nil {
				logger.Errorf("failed to get %s rate for provider %s: %v", orderToken.Edges.Token.Symbol, providerID, err)
				continue
			}

			// Check provider's rate against the market rate to ensure it's not too far off
			percentDeviation := utils.AbsPercentageDeviation(bucket.Edges.Currency.MarketRate, rate)

			if serverConf.Environment == "production" && percentDeviation.GreaterThan(orderConf.PercentDeviationFromMarketRate) {
				// Skip this provider if the rate is too far off
				// TODO: add a logic to notify the provider(s) to update his rate since it's stale. could be a cron job
				continue
			}

			// Serialize the provider ID, token, rate, min and max order amount into a single string
			data := fmt.Sprintf("%s:%s:%s:%s:%s", providerID, orderToken.Edges.Token.Symbol, rate, orderToken.MinOrderAmount, orderToken.MaxOrderAmount)

			// Enqueue the serialized data into the circular queue
			err = storage.RedisClient.RPush(ctx, redisKey, data).Err()
			if err != nil {
				logger.Errorf("failed to enqueue provider data to circular queue: %v", err)
			}
		}
	}
}

// AssignLockPaymentOrders assigns lock payment orders to providers
func (s *PriorityQueueService) AssignLockPaymentOrder(ctx context.Context, order types.LockPaymentOrderFields) error {
	orderIDPrefix := strings.Split(order.ID.String(), "-")[0]

	excludeList, err := storage.RedisClient.LRange(ctx, fmt.Sprintf("order_exclude_list_%s", order.ID), 0, -1).Result()
	if err != nil {
		logger.Errorf("%s - failed to get exclude list: %v", order.ID, err)
		return err
	}

	// Sends order directly to the specified provider in order.
	// Incase of failure, do nothing. The order will eventually refund
	if order.ProviderID != "" && !utils.ContainsString(excludeList, order.ProviderID) {
		provider, err := storage.Client.ProviderProfile.
			Query().
			Where(
				providerprofile.IDEQ(order.ProviderID),
			).
			Only(ctx)

		if err == nil {
			// TODO: check for provider's minimum and maximum rate for negotiation
			// Update the rate with the current rate if order was last updated more than 10 mins ago
			if order.UpdatedAt.Before(time.Now().Add(-10 * time.Minute)) {
				order.Rate, err = s.GetProviderRate(ctx, provider, order.Token.Symbol, order.ProvisionBucket.Edges.Currency.Code)
				if err != nil {
					logger.Errorf("%s - failed to get rate for provider %s: %v", orderIDPrefix, order.ProviderID, err)
				}
				_, err = storage.Client.PaymentOrder.
					Update().
					Where(paymentorder.IDEQ(order.ID)).
					SetRate(order.Rate).
					Save(ctx)
				if err != nil {
					logger.Errorf("%s - failed to update rate for provider %s: %v", orderIDPrefix, order.ProviderID, err)
				}
			}
			err = s.sendOrderRequest(ctx, order)
			if err == nil {
				return nil
			}
			logger.Errorf("%s - failed to send order request to specific provider %s: %v", orderIDPrefix, order.ProviderID, err)
		} else {
			logger.Errorf("%s - failed to get provider: %v", orderIDPrefix, err)
		}

		if provider.VisibilityMode == providerprofile.VisibilityModePrivate {
			return nil
		}
	}

	// Get the first provider from the circular queue
	redisKey := fmt.Sprintf("bucket_%s_%s_%s", order.ProvisionBucket.Edges.Currency.Code, order.ProvisionBucket.MinAmount, order.ProvisionBucket.MaxAmount)

	// partnerProviders := []string{}

	for index := 0; ; index++ {
		providerData, err := storage.RedisClient.LIndex(ctx, redisKey, int64(index)).Result()
		if err != nil {
			break
		}

		// if providerData == "" {
		// 	// Reached the end of the queue
		// 	logger.Errorf("%s - rate didn't match a provider, finding a partner provider", orderIDPrefix)

		// 	if len(partnerProviders) == 0 {
		// 		logger.Errorf("%s - no partner providers found", orderIDPrefix)
		// 		return nil
		// 	}

		// 	// Pick a random partner provider
		// 	randomIndex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(partnerProviders))
		// 	providerData = partnerProviders[randomIndex]
		// }

		// Extract the rate from the data (assuming it's in the format "providerID:token:rate:minAmount:maxAmount")
		parts := strings.Split(providerData, ":")
		if len(parts) != 5 {
			logger.Errorf("%s - invalid data format at index %d: %s", orderIDPrefix, index, providerData)
			continue // Skip this entry due to invalid format
		}

		order.ProviderID = parts[0]

		// Skip entry if provider is excluded
		if utils.ContainsString(excludeList, order.ProviderID) {
			continue
		}

		// Skip entry if token doesn't match
		if parts[1] != order.Token.Symbol {
			continue
		}

		// Skip entry if order amount is not within provider's min and max order amount
		minOrderAmount, err := decimal.NewFromString(parts[3])
		if err != nil {
			continue
		}

		maxOrderAmount, err := decimal.NewFromString(parts[4])
		if err != nil {
			continue
		}

		if order.Amount.LessThan(minOrderAmount) || order.Amount.GreaterThan(maxOrderAmount) {
			continue
		}

		// Fetch and check provider for rate match
		rate, err := decimal.NewFromString(parts[2])
		if err != nil {
			continue
		}

		// TODO: make the slippage of 0.5 configurable by provider
		if rate.Sub(order.Rate).Abs().LessThanOrEqual(decimal.NewFromFloat(0.5)) {
			// Found a match for the rate
			if index == 0 {
				// Match found at index 0, perform LPOP to dequeue
				data, err := storage.RedisClient.LPop(ctx, redisKey).Result()
				if err != nil {
					logger.Errorf("%s - failed to dequeue from circular queue: %v", orderIDPrefix, err)
					return err
				}

				// Enqueue data to the end of the queue
				err = storage.RedisClient.RPush(ctx, redisKey, data).Err()
				if err != nil {
					logger.Errorf("%s - failed to enqueue to circular queue: %v", orderIDPrefix, err)
					return err
				}
			}

			// Assign the order to the provider and save it to Redis
			err = s.sendOrderRequest(ctx, order)
			if err != nil {
				logger.Errorf("%s - failed to send order request to specific provider %s: %v", orderIDPrefix, order.ProviderID, err)

				// Push provider ID to order exclude list
				orderKey := fmt.Sprintf("order_exclude_list_%s", order.ID)
				_, err = storage.RedisClient.RPush(ctx, orderKey, order.ProviderID).Result()
				if err != nil {
					logger.Errorf("%s - error pushing provider %s to order_exclude_list on Redis: %v", orderIDPrefix, order.ProviderID, err)
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

	// TODO: Now we need to add currency
	orderRequestData := map[string]interface{}{
		"amount":      order.Amount.Mul(order.Rate).RoundBank(0).String(),
		"institution": order.Institution,
		"providerId":  order.ProviderID,
	}

	if err := storage.RedisClient.HSet(ctx, orderKey, orderRequestData).Err(); err != nil {
		logger.Errorf("failed to map order to a provider in Redis: %v", err)
		return err
	}

	// Set a TTL for the order request
	err := storage.RedisClient.ExpireAt(ctx, orderKey, time.Now().Add(orderConf.OrderRequestValidity)).Err()
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
		Where(providerprofile.IDEQ(providerID)).
		WithAPIKey().
		Only(ctx)
	if err != nil {
		return err
	}

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
	_, err = fastshot.NewClient(provider.HostIdentifier).
		Config().SetTimeout(30*time.Second).
		Header().Add("X-Request-Signature", signature).
		Build().POST("/new_order").
		Body().AsJSON(orderRequestData).
		Send()
	if err != nil {
		return err
	}

	return nil
}
