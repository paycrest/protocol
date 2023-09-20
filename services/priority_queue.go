package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/provideravailability"
	"github.com/paycrest/paycrest-protocol/ent/providerprofile"
	"github.com/paycrest/paycrest-protocol/ent/providerrating"
	"github.com/paycrest/paycrest-protocol/ent/provisionbucket"
	"github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/redis/go-redis/v9"
)

type PriorityQueueService struct{}

// NewPriorityQueueService creates a new instance of PriorityQueueService
func NewPriorityQueueService() *PriorityQueueService {
	return &PriorityQueueService{}
}

// ProcessBucketQueues creates a priority queue for each bucket and saves it to redis
func (s *PriorityQueueService) ProcessBucketQueues(ctx context.Context) error {

	buckets, err := s.GetProvidersByBucket(ctx)
	if err != nil {
		return fmt.Errorf("failed to process bucket queues: %w", err)
	}

	var wg sync.WaitGroup

	for _, bucket := range buckets {
		wg.Add(1)
		go s.CreatePriorityQueueForBucket(ctx, bucket)
	}

	wg.Wait()

	return nil
}

// GetProvidersByBucket returns a list of providers grouped by bucket
func (s *PriorityQueueService) GetProvidersByBucket(ctx context.Context) ([]*ent.ProvisionBucket, error) {
	buckets, err := storage.Client.ProvisionBucket.
		Query().
		Select(provisionbucket.EdgeProviderProfiles).
		WithProviderProfiles(func(ppq *ent.ProviderProfileQuery) {
			ppq.WithProviderRating(func(prq *ent.ProviderRatingQuery) {
				prq.Select(providerrating.FieldTrustScore)
			})
			ppq.Select(providerprofile.FieldID)

			// Filter only providers that are always available
			// or are available until one hour from now
			// TODO: the duration should be a config setting
			oneHourFromNow := time.Now().Add(time.Hour)
			ppq.Where(
				providerprofile.HasAvailabilityWith(
					provideravailability.And(
						provideravailability.CadenceEQ(provideravailability.CadenceAlways),
						provideravailability.EndTimeGTE(oneHourFromNow),
					),
				),
			)
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return buckets, nil
}

// CreatePriorityQueueForBucket creates a priority queue for a bucket and saves it to redis
func (s *PriorityQueueService) CreatePriorityQueueForBucket(ctx context.Context, bucket *ent.ProvisionBucket) {
	// Create a slice to store the sorted set members with their scores
	providers := bucket.Edges.ProviderProfiles
	members := make([]redis.Z, len(providers))

	// Populate the members slice with providers and their trust scores
	for i, provider := range providers {
		trustScore, _ := provider.Edges.ProviderRating.TrustScore.Float64()

		members[i] = redis.Z{
			Score:  trustScore,
			Member: provider.ID,
		}
	}

	// Add bucket with sorted priority queue to the redis cache
	// e.g {"bucket_<id>": [1,2,3,4,5,6,7]}
	redisKey := fmt.Sprintf("bucket_%d", bucket.ID)

	// Add the members to the sorted set
	_, err := storage.RedisClient.ZAdd(ctx, redisKey, members...).Result()
	if err != nil {
		logger.Errorf("failed to add bucket priority queue to Redis: %v", err)
	}
}

// AssignLockPaymentOrders assigns lock payment orders to providers
func (s *PriorityQueueService) AssignLockPaymentOrder(ctx context.Context, order LockPaymentOrderFields) error {
	// Get the first provider from the priority queue
	redisKey := fmt.Sprintf("bucket_%d", order.ProvisionBucket.ID)
	providerIDs, err := storage.RedisClient.ZRevRange(ctx, redisKey, 0, 0).Result()
	if err != nil {
		logger.Errorf("failed to get provider from priority queue: %v", err)
		return err
	}

	if len(providerIDs) == 0 {
		logger.Errorf("no providers available for bucket %d", order.ProvisionBucket.ID)
		return fmt.Errorf("no providers available for bucket %d", order.ProvisionBucket.ID)
	}

	// Assign the order to the provider and save it to redis
	orderKey := fmt.Sprintf("order_request_%d", order.ID)
	data := map[string]interface{}{
		"amount":      order.Amount.Mul(order.Rate),
		"token":       order.Token.Symbol,
		"institution": order.Institution,
		"provider_id": providerIDs[0],
	}

	_, err = storage.RedisClient.HSet(ctx, orderKey, data).Result()
	if err != nil {
		logger.Errorf("failed to map order to a provider in redis: %v", err)
		return err
	}

	// Set a TTL for the order request
	_, err = storage.RedisClient.Expire(ctx, orderKey, OrderConf.OrderRequestValidity).Result()
	if err != nil {
		logger.Errorf("failed to set TTL for order request: %v", err)
		return err
	}

	// Remove the provider from the priority queue
	_, err = storage.RedisClient.ZRem(ctx, redisKey, providerIDs[0]).Result()
	if err != nil {
		logger.Errorf("failed to remove provider from priority queue: %v", err)
		return err
	}

	// Create a priority queue for the bucket if there was only one provider
	if len(providerIDs) == 1 {
		s.CreatePriorityQueueForBucket(ctx, order.ProvisionBucket)
	}

	// TODO: Send wss message to the provider (for automatic provider case)

	// TODO: Send out a push notification to the provider (for manual provider case)

	return nil
}
