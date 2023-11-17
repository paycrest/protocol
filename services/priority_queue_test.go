package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/ent/migrate"
	"github.com/paycrest/protocol/storage"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/utils/test"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestProviderVisibilityMode(t *testing.T) {

	// Set up test Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()
	db.Client = client

	// Run the auto migration tool.
	err := client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true))
	assert.NoError(t, err)

	// Seed the database
	err = storage.SeedAll()
	assert.NoError(t, err)

	// Initialize Redis
	err = storage.InitializeRedis()
	assert.NoError(t, err)

	// Set up test service
	pqService := NewPriorityQueueService()

	// Setup test data
	err = setup()
	assert.NoError(t, err)

	// Set up test order
	// order := types.LockPaymentOrderFields{
	// 	ID: 1,
	// 	ProvisionBucket: &types.Provi{
	// 		Edges: types.ProvisionBucketEdges{
	// 			Currency: &types.Currency{
	// 				Code: "USD",
	// 			},
	// 		},
	// 		MinAmount: decimal.NewFromFloat(100),
	// 		MaxAmount: decimal.NewFromFloat(1000),
	// 	},
	// 	Rate:        decimal.NewFromFloat(1.5),
	// 	Amount:      decimal.NewFromFloat(500),
	// 	Token:       &types.Token{Symbol: "USDT"},
	// 	Institution: "Test Institution",
	// 	ProviderID:  "",
	// }

	// Set up test provider user
	ProviderPublic, err := test.CreateTestUser(map[string]string{
		"scope": "provider",
		"email": "public@test.com",
	})
	assert.NoError(t, err)

	providerPrivate, err := test.CreateTestUser(map[string]string{
		"scope": "provider",
		"email": "private@test.com",
	})
	assert.NoError(t, err)

	// Set up test provider currency
	currency, err := test.CreateTestFiatCurrency(nil)
	assert.NoError(t, err)

	_, err = test.CreateTestProviderProfile(nil, ProviderPublic, currency)
	assert.NoError(t, err)
	_, err = test.CreateTestProviderProfile(nil, providerPrivate, currency)
	assert.NoError(t, err)

	t.Run("GetProvidersByBucket", func(t *testing.T) {
		t.Run("Return only providers with visibility mode public", func(t *testing.T) {
			buckets, err := pqService.GetProvidersByBucket(context.Background())
			assert.NoError(t, err)
			fmt.Printf("!!!buckets >> %v", buckets)

		})
	})

	// service.GetProvidersByBucket(ctx context.Context)

	// Test case 1: no exclude list, no specified provider, circular queue has a match
	// err := redisClient.Del(ctx, "bucket_USD_100_1000").Err()
	// assert.NoError(t, err)

	// err = redisClient.RPush(ctx, "bucket_USD_100_1000", "p1:1.5", "p2:1.6", "p3:1.4").Err()
	// assert.NoError(t, err)

	// err = service.AssignLockPaymentOrder(ctx, order)
	// assert.NoError(t, err)

	// orderKey := fmt.Sprintf("order_request_%d", order.ID)
	// orderRequestData, err := redisClient.HGetAll(ctx, orderKey).Result()
	// assert.NoError(t, err)

	// assert.Equal(t, "750", orderRequestData["amount"])
	// assert.Equal(t, "USDT", orderRequestData["token"])
	// assert.Equal(t, "Test Institution", orderRequestData["institution"])
	// assert.Equal(t, "p1", orderRequestData["providerId"])

	// Test case 2: no exclude list, no specified provider, circular queue has no match
	// 	err = redisClient.Del(ctx, "bucket_USD_100_1000").Err()
	// 	assert.NoError(t, err)

	// 	err = redisClient.RPush(ctx, "bucket_USD_100_1000", "p1:1.6", "p2:1.7", "p3:1.8").Err()
	// 	assert.NoError(t, err)

	// 	err = service.AssignLockPaymentOrder(ctx, order)
	// 	assert.NoError(t, err)

	// 	orderRequestData, err = redisClient.HGetAll(ctx, orderKey).Result()
	// 	assert.NoError(t, err)

	// 	assert.Equal(t, "750", orderRequestData["amount"])
	// 	assert.Equal(t, "USDT", orderRequestData["token"])
	// 	assert.Equal(t, "Test Institution", orderRequestData["institution"])
	// 	assert.Equal(t, "p1", orderRequestData["providerId"])

	// 	// Test case 3: no exclude list, specified provider has no match
	// 	err = redisClient.Del(ctx, "bucket_USD_100_1000").Err()
	// 	assert.NoError(t, err)

	// 	err = redisClient.RPush(ctx, "bucket_USD_100_1000", "p1:1.6", "p2:1.7", "p3:1.8").Err()
	// 	assert.NoError(t, err)

	// 	order.ProviderID = "p4"

	// 	err = service.AssignLockPaymentOrder(ctx, order)
	// 	assert.NoError(t, err)

	// 	orderRequestData, err = redisClient.HGetAll(ctx, orderKey).Result()
	// 	assert.NoError(t, err)

	// 	assert.Equal(t, "750", orderRequestData["amount"])
	// 	assert.Equal(t, "USDT", orderRequestData["token"])
	// 	assert.Equal(t, "Test Institution", orderRequestData["institution"])
	// 	assert.Equal(t, "p4", orderRequestData["providerId"])

	// 	// Test case 4: exclude list, specified provider has no match
	// 	err = redisClient.Del(ctx, "bucket_USD_100_1000").Err()
	// 	assert.NoError(t, err)

	// 	err = redisClient.RPush(ctx, "bucket_USD_100_1000", "p1:1.6", "p2:1.7", "p3:1.8").Err()
	// 	assert.NoError(t, err)

	// 	err = redisClient.RPush(ctx, "order_exclude_list_1", "p4").Err()
	// 	assert.NoError(t, err)

	// 	order.ProviderID = "p4"

	// 	err = service.AssignLockPaymentOrder(ctx, order)
	// 	assert.NoError(t, err)

	// 	orderRequestData, err = redisClient.HGetAll(ctx, orderKey).Result()
	// 	assert.NoError(t, err)

	// 	assert.Equal(t, "750", orderRequestData["amount"])
	// 	assert.Equal(t, "USDT", orderRequestData["token"])
	// 	assert.Equal(t, "Test Institution", orderRequestData["institution"])
	// 	assert.Equal(t, "p1", orderRequestData["providerId"])

	// 	// Test case 5: exclude list, specified provider has a match
	// 	err = redisClient.Del(ctx, "bucket_USD_100_1000").Err()
	// 	assert.NoError(t, err)

	// 	err = redisClient.RPush(ctx, "bucket_USD_100_1000", "p1:1.5", "p2:1.6", "p3:1.4").Err()
	// 	assert.NoError(t, err)

	// 	err = redisClient.RPush(ctx, "order_exclude_list_1", "p2").Err()
	// 	assert.NoError(t, err)

	// 	order.ProviderID = "p2"

	// 	err = service.AssignLockPaymentOrder(ctx, order)
	// 	assert.NoError(t, err)

	// 	orderRequestData, err = redisClient.HGetAll(ctx, orderKey).Result()
	// 	assert.NoError(t, err)

	// 	assert.Equal(t, "750", orderRequestData["amount"])
	// 	assert.Equal(t, "USDT", orderRequestData["token"])
	// 	assert.Equal(t, "Test Institution", orderRequestData["institution"])
	// 	assert.Equal(t, "p1", orderRequestData["providerId"])

	// 	// Test case 6: invalid provider data format
	// 	err = redisClient.Del(ctx, "bucket_USD_100_1000").Err()
	// 	assert.NoError(t, err)

	// 	err = redisClient.RPush(ctx, "bucket_USD_100_1000", "p1:1.5", "invalid_data", "p3:1.4").Err()
	// 	assert.NoError(t, err)

	// 	err = service.AssignLockPaymentOrder(ctx, order)
	// 	assert.NoError(t, err)

	// 	orderRequestData, err = redisClient.HGetAll(ctx, orderKey).Result()
	// 	assert.NoError(t, err)

	// 	assert.Equal(t, "750", orderRequestData["amount"])
	// 	assert.Equal(t, "USDT", orderRequestData["token"])
	// 	assert.Equal(t, "Test Institution", orderRequestData["institution"])
	// 	assert.Equal(t, "p1", orderRequestData["providerId"])

	// 	// Test case 7: Redis error
	// 	err = redisClient.Del(ctx, "bucket_USD_100_1000").Err()
	// 	assert.NoError(t, err)

	// 	err = redisClient.RPush(ctx, "bucket_USD_100_1000", "p1:1.5", "p2:1.6", "p3:1.4").Err()
	// 	assert.NoError(t, err)

	// 	// Force a Redis error by deleting the Redis key
	// 	err = redisClient.Del(ctx, orderKey).Err()
	// 	assert.NoError(t, err)

	// 	err = service.AssignLockPaymentOrder(ctx, order)
	// 	assert.Error(t, err)
	// 	assert.True(t, strings.HasPrefix(err.Error(), "failed to map order to a provider in Redis: "))

	// 	// Test case 8: provider notification error
	// 	err = redisClient.Del(ctx, "bucket_USD_100_1000").Err()
	// 	assert.NoError(t, err)

	// 	err = redisClient.RPush(ctx, "bucket_USD_100_1000", "p1:1.5", "p2:1.6", "p3:1.4").Err()
	// 	assert.NoError(t, err)

	// 	// Force a provider notification error by using an invalid provider ID
	// 	order.ProviderID = "invalid_provider_id"

	// 	err = service.AssignLockPaymentOrder(ctx, order)
	// 	assert.NoError(t, err)

	// 	orderRequestData, err = redisClient.HGetAll(ctx, orderKey).Result()
	// 	assert.NoError(t, err)

	// 	assert.Equal(t, "750", orderRequestData["amount"])
	// 	assert.Equal(t, "USDT", orderRequestData["token"])
	// 	assert.Equal(t, "Test Institution", orderRequestData["institution"])
	// 	assert.Equal(t, "p1", orderRequestData["providerId"])

}

func NewProfileController() {
	panic("unimplemented")
}
