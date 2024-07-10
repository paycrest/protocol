package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/provisionbucket"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	cryptoUtils "github.com/paycrest/protocol/utils/crypto"
	"github.com/paycrest/protocol/utils/test"
	tokenUtils "github.com/paycrest/protocol/utils/token"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var testCtxForPQ = struct {
	user                     *ent.User
	providerProfile          *ent.ProviderProfile
	providerProfileAPIsecret string

	privateProviderPrivate *ent.ProviderProfile
	currency               *ent.FiatCurrency
	client                 types.RPCClient
	token                  *ent.Token
	minAmount              decimal.Decimal
	maxAmount              decimal.Decimal
	bucket                 *ent.ProvisionBucket
}{}

func setupForPQ() error {
	// Set up test data
	testCtxForPQ.maxAmount = decimal.NewFromFloat(10000)
	testCtxForPQ.minAmount = decimal.NewFromFloat(1)

	backend, err := test.SetUpTestBlockchain()
	if err != nil {
		return err
	}
	token, err := test.CreateERC20Token(backend, map[string]interface{}{})
	if err != nil {
		return err
	}
	testCtxForPQ.token = token

	user, err := test.CreateTestUser(map[string]interface{}{
		"scope": "provider",
		"email": "providerjohndoe@test.com",
	})
	if err != nil {
		return err
	}
	testCtxForPQ.user = user

	currency, err := test.CreateTestFiatCurrency(map[string]interface{}{
		"code":        "KES",
		"short_name":  "Shilling",
		"decimals":    2,
		"symbol":      "KSh",
		"name":        "Kenyan Shilling",
		"market_rate": 550.0,
	})
	if err != nil {
		return err
	}
	testCtxForPQ.currency = currency

	providerProfile, err := test.CreateTestProviderProfile(map[string]interface{}{
		"user_id":         testCtxForPQ.user.ID,
		"currency_id":     currency.ID,
		"host_identifier": "https://example2.com",
	})
	if err != nil {
		return err
	}
	apiKeyService := NewAPIKeyService()
	secret, _, err := apiKeyService.GenerateAPIKey(
		context.Background(),
		nil,
		nil,
		providerProfile,
	)
	if err != nil {
		return err
	}
	testCtxForPQ.providerProfileAPIsecret = secret.Secret
	_, err = test.AddProviderOrderTokenToProvider(
		map[string]interface{}{
			"fixed_conversion_rate":    decimal.NewFromFloat(100),
			"conversion_rate_type":     "fixed",
			"floating_conversion_rate": decimal.NewFromFloat(1.0),
			"max_order_amount":         decimal.NewFromFloat(1000),
			"min_order_amount":         decimal.NewFromFloat(1.0),
			"tokenSymbol":              token.Symbol,
			"provider":                 providerProfile,
		},
	)
	if err != nil {
		return err
	}
	testCtxForPQ.providerProfile = providerProfile

	bucket, err := test.CreateTestProvisionBucket(map[string]interface{}{
		"provider_id": providerProfile.ID,
		"min_amount":  decimal.NewFromFloat(1),
		"max_amount":  decimal.NewFromFloat(10000.0),
		"currency_id": currency.ID,
	})
	if err != nil {
		return err
	}
	testCtxForPQ.bucket = bucket

	providerPrivate, err := test.CreateTestUser(map[string]interface{}{
		"scope": "provider",
		"email": "private@test.com",
	})
	if err != nil {
		return err
	}

	privateProviderPrivate, err := test.CreateTestProviderProfile(map[string]interface{}{
		"currency_id":     currency.ID,
		"visibility_mode": "private",
		"user_id":         providerPrivate.ID,
	})
	if err != nil {
		return err
	}
	testCtxForPQ.privateProviderPrivate = privateProviderPrivate

	_, err = test.CreateTestProvisionBucket(map[string]interface{}{
		"provider_id": privateProviderPrivate.ID,
		"min_amount":  testCtxForPQ.minAmount,
		"max_amount":  testCtxForPQ.maxAmount,
		"currency_id": currency.ID,
	})
	if err != nil {
		return err
	}

	// Set up payment order
	_, err = test.CreateTestLockPaymentOrder(map[string]interface{}{
		"provider": privateProviderPrivate,
		"tokenID":  testCtxForPQ.token.ID})
	if err != nil {
		return err
	}
	_, err = test.CreateTestLockPaymentOrder(map[string]interface{}{
		"provider": providerProfile,
		"tokenID":  testCtxForPQ.token.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
func TestPriorityQueueTest(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer redisClient.Close()

	db.RedisClient = redisClient
	{

		err := redisClient.FlushAll(context.Background()).Err()
		assert.NoError(t, err)
	}

	db.Client = client

	// Setup test data
	err := setupForPQ()
	assert.NoError(t, err)

	service := NewPriorityQueueService()
	t.Run("TestGetProvisionBuckets", func(t *testing.T) {
		buckets, err := service.GetProvisionBuckets(context.Background())
		assert.NoError(t, err)
		assert.Greater(t, len(buckets), 0)

	})

	t.Run("TestCreatePriorityQueueForBucket", func(t *testing.T) {
		ctx := context.Background()
		bucket, err := test.CreateTestProvisionBucket(map[string]interface{}{
			"provider_id": testCtxForPQ.privateProviderPrivate.ID,
			"min_amount":  testCtxForPQ.minAmount,
			"max_amount":  testCtxForPQ.maxAmount,
			"currency_id": testCtxForPQ.currency.ID,
		})
		assert.NoError(t, err)

		_bucket, err := db.Client.ProvisionBucket.
			Query().
			Where(provisionbucket.IDEQ(bucket.ID)).
			WithCurrency().
			WithProviderProfiles().
			Only(ctx)
		assert.NoError(t, err)

		service.CreatePriorityQueueForBucket(ctx, _bucket)

		redisKey := fmt.Sprintf("bucket_%s_%s_%s", _bucket.Edges.Currency.Code, testCtxForPQ.minAmount, testCtxForPQ.maxAmount)

		// err = db.RedisClient.RPush(ctx, redisKey, testCtxForPQ.privateProviderPrivate.ID).Err()
		// assert.NoError(t, err)
		data, err := db.RedisClient.LRange(ctx, redisKey, 0, -1).Result()
		assert.NoError(t, err)
		assert.Equal(t, len(data), 1)
		assert.Contains(t, data[0], testCtxForPQ.privateProviderPrivate.ID)

	})
	t.Run("TestProcessBucketQueues", func(t *testing.T) {
		err = service.ProcessBucketQueues()
		assert.NoError(t, err)

		redisKey := fmt.Sprintf("bucket_%s_%s_%s", testCtxForPQ.currency.Code, testCtxForPQ.minAmount, testCtxForPQ.maxAmount)

		data, err := db.RedisClient.LRange(context.Background(), redisKey, 0, -1).Result()
		assert.NoError(t, err)
		assert.Equal(t, len(data), 1)
	})

	t.Run("TestAssignLockPaymentOrder", func(t *testing.T) {

		bucket, err := test.CreateTestProvisionBucket(map[string]interface{}{
			"provider_id": testCtxForPQ.privateProviderPrivate.ID,
			"min_amount":  testCtxForPQ.minAmount,
			"max_amount":  testCtxForPQ.maxAmount,
			"currency_id": testCtxForPQ.currency.ID,
		})
		assert.NoError(t, err)
		_order, err := test.CreateTestLockPaymentOrder(map[string]interface{}{"provider": testCtxForPQ.providerProfile, "tokenID": testCtxForPQ.token.ID})
		assert.NoError(t, err)

		_, err = test.AddProvisionBucketToLockPaymentOrder(_order, bucket.ID)
		assert.NoError(t, err)

		err = db.RedisClient.RPush(context.Background(), fmt.Sprintf("order_exclude_list_%s", _order.ID), testCtxForPQ.providerProfile.ID).Err()
		assert.NoError(t, err)

		order, err := db.Client.LockPaymentOrder.
			Query().
			Where(lockpaymentorder.IDEQ(_order.ID)).
			WithProvisionBucket(func(pb *ent.ProvisionBucketQuery) {
				pb.WithCurrency()
			}).
			WithToken().
			Only(context.Background())

		assert.NoError(t, err)

		err = service.AssignLockPaymentOrder(context.Background(), types.LockPaymentOrderFields{
			ID:                order.ID,
			Token:             testCtxForPQ.token,
			GatewayID:         order.GatewayID,
			Amount:            order.Amount,
			Rate:              order.Rate,
			BlockNumber:       order.BlockNumber,
			Institution:       order.Institution,
			AccountIdentifier: order.AccountIdentifier,
			AccountName:       order.AccountName,
			Memo:              order.Memo,
			ProvisionBucket:   order.Edges.ProvisionBucket,
		})
		assert.NoError(t, err)
	})

	t.Run("TestGetProviderRate", func(t *testing.T) {
		rate, err := service.GetProviderRate(context.Background(), testCtxForPQ.providerProfile)
		assert.NoError(t, err)
		_rate, ok := rate.Float64()
		assert.True(t, ok)
		assert.Equal(t, _rate, float64(100))
	})

	t.Run("TestSendOrderRequest", func(t *testing.T) {
		bucket, err := test.CreateTestProvisionBucket(map[string]interface{}{
			"provider_id": testCtxForPQ.privateProviderPrivate.ID,
			"min_amount":  testCtxForPQ.minAmount,
			"max_amount":  testCtxForPQ.maxAmount,
			"currency_id": testCtxForPQ.currency.ID,
		})
		assert.NoError(t, err)
		_order, err := test.CreateTestLockPaymentOrder(map[string]interface{}{"provider": testCtxForPQ.providerProfile, "tokenID": testCtxForPQ.token.ID})
		assert.NoError(t, err)

		_, err = test.AddProvisionBucketToLockPaymentOrder(_order, bucket.ID)
		assert.NoError(t, err)

		err = db.RedisClient.RPush(context.Background(), fmt.Sprintf("order_exclude_list_%s", _order.ID), testCtxForPQ.providerProfile.ID).Err()
		assert.NoError(t, err)

		order, err := db.Client.LockPaymentOrder.
			Query().
			Where(lockpaymentorder.IDEQ(_order.ID)).
			WithProvisionBucket(func(pb *ent.ProvisionBucketQuery) {
				pb.WithCurrency()
			}).
			WithToken().
			Only(context.Background())

		assert.NoError(t, err)

		err = service.sendOrderRequest(context.Background(), types.LockPaymentOrderFields{
			ID:                order.ID,
			ProviderID:        testCtxForPQ.providerProfile.ID,
			Token:             testCtxForPQ.token,
			GatewayID:         order.GatewayID,
			Amount:            order.Amount,
			Rate:              order.Rate,
			BlockNumber:       order.BlockNumber,
			Institution:       order.Institution,
			AccountIdentifier: order.AccountIdentifier,
			AccountName:       order.AccountName,
			Memo:              order.Memo,
			ProvisionBucket:   order.Edges.ProvisionBucket,
		})
		assert.NoError(t, err)
		t.Run("TestNotifyProvider", func(t *testing.T) {

			// setup httpmock
			httpmock.Activate()
			defer httpmock.Deactivate()

			httpmock.RegisterResponder("POST", testCtxForPQ.providerProfile.HostIdentifier+"/new_order",
				func(r *http.Request) (*http.Response, error) {
					bytes, err := io.ReadAll(r.Body)
					if err != nil {
						log.Fatal(err)
					}
					// Compute HMAC
					decodedSecret, err := base64.StdEncoding.DecodeString(testCtxForPQ.providerProfileAPIsecret)
					assert.NoError(t, err)
					decryptedSecret, err := cryptoUtils.DecryptPlain(decodedSecret)
					assert.NoError(t, err)
					signature := tokenUtils.GenerateHMACSignature(map[string]interface{}{
						"data": "test",
					}, string(decryptedSecret))
					assert.Equal(t, r.Header.Get("X-Request-Signature"), signature)
					if strings.Contains(string(bytes), "data") && strings.Contains(string(bytes), "test") {
						resp := httpmock.NewBytesResponse(200, nil)
						return resp, nil
					} else {
						return nil, nil
					}
				},
			)
			err := service.notifyProvider(context.Background(), map[string]interface{}{
				"providerId": testCtxForPQ.providerProfile.ID,
				"data":       "test",
			})
			assert.NoError(t, err)
		})
	})

	t.Run("TestNoErrorFunctions", func(t *testing.T) {

		t.Run("TestReassignUnfulfilledLockOrders", func(t *testing.T) {
			// redisKey := fmt.Sprintf("bucket_%s_%s_%s", testCtxForPQ.currency.Code, testCtxForPQ.minAmount, testCtxForPQ.maxAmount)

			// pubsub := db.RedisClient.Subscribe(context.Background(), redisKey)

			// // Listen for messages
			// go func() {
			// 	for msg := range pubsub.Channel() {
			// 		assert.Equal(t, msg, "")
			// 	}
			// }()

			service.ReassignUnfulfilledLockOrders()
			// Keep the main thread alive
			// select {}
		})
		t.Run("TestReassignStaleOrderRequest", func(t *testing.T) {
		})

		t.Run("TestReassignUnValidatedLockOrders", func(t *testing.T) {
		})
		t.Run("TestReassignPendingOrders", func(t *testing.T) {
		})
	})
}

// func TestProviderVisibilityMode(t *testing.T) {

// 	// Set up test Redis client
// 	redisClient := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 	})
// 	defer redisClient.Close()

// 	// Set up test database client
// 	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
// 	defer client.Close()
// 	db.Client = client

// 	// Run the auto migration tool.
// 	err := client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true))
// 	assert.NoError(t, err)

// 	// Seed the database
// 	err = storage.SeedAll()
// 	assert.NoError(t, err)

// 	// Initialize Redis
// 	err = storage.InitializeRedis()
// 	assert.NoError(t, err)

// 	// Set up test service
// 	pqService := NewPriorityQueueService()

// 	// Set up test order
// 	// order := types.LockPaymentOrderFields{
// 	// 	ID: 1,
// 	// 	ProvisionBucket: &types.Provi{
// 	// 		Edges: types.ProvisionBucketEdges{
// 	// 			Currency: &types.Currency{
// 	// 				Code: "USD",
// 	// 			},
// 	// 		},
// 	// 		MinAmount: decimal.NewFromFloat(100),
// 	// 		MaxAmount: decimal.NewFromFloat(1000),
// 	// 	},
// 	// 	Rate:        decimal.NewFromFloat(1.5),
// 	// 	Amount:      decimal.NewFromFloat(500),
// 	// 	Token:       &types.Token{Symbol: "USDT"},
// 	// 	Institution: "Test Institution",
// 	// 	ProviderID:  "",
// 	// }

// 	// Set up test provider user
// 	ProviderPublic, err := test.CreateTestUser(map[string]string{
// 		"scope": "provider",
// 		"email": "public@test.com",
// 	})
// 	assert.NoError(t, err)

// 	providerPrivate, err := test.CreateTestUser(map[string]string{
// 		"scope": "provider",
// 		"email": "private@test.com",
// 	})
// 	assert.NoError(t, err)

// 	// Set up test provider currency
// 	currency, err := test.CreateTestFiatCurrency(nil)
// 	assert.NoError(t, err)

// 	publicProviderPrivate, err := test.CreateTestProviderProfile(nil, ProviderPublic, currency)
// 	assert.NoError(t, err)
// 	privateProviderPrivate, err := test.CreateTestProviderProfile(map[string]interface{}{"visibility_mode": "private"}, providerPrivate, currency)
// 	assert.NoError(t, err)

// 	// Set up payment order
// 	_, err = test.CreateTestLockPaymentOrder(map[string]interface{}{"provider": privateProviderPrivate})
// 	assert.NoError(t, err)
// 	_, err = test.CreateTestLockPaymentOrder(map[string]interface{}{"provider": publicProviderPrivate, "order_id": "order_1234"})
// 	assert.NoError(t, err)

// 	t.Run("GetProvidersByBucket", func(t *testing.T) {
// 		t.Run("Return only providers with visibility mode public", func(t *testing.T) {
// 			buckets, err := pqService.GetProvidersByBucket(context.Background())
// 			assert.NoError(t, err)
// 			fmt.Printf("!!!buckets >> %v", buckets)

// 		})
// 	})

// 	service.GetProvidersByBucket(context.Background())

// 	// Test case 1: no exclude list, no specified provider, circular queue has a match
// 	err := redisClient.Del(ctx, "bucket_USD_100_1000").Err()
// 	assert.NoError(t, err)

// 	err = redisClient.RPush(ctx, "bucket_USD_100_1000", "p1:1.5", "p2:1.6", "p3:1.4").Err()
// 	assert.NoError(t, err)

// 	err = service.AssignLockPaymentOrder(ctx, order)
// 	assert.NoError(t, err)

// 	orderKey := fmt.Sprintf("order_request_%d", order.ID)
// 	orderRequestData, err := redisClient.HGetAll(ctx, orderKey).Result()
// 	assert.NoError(t, err)

// 	assert.Equal(t, "750", orderRequestData["amount"])
// 	assert.Equal(t, "USDT", orderRequestData["token"])
// 	assert.Equal(t, "Test Institution", orderRequestData["institution"])
// 	assert.Equal(t, "p1", orderRequestData["providerId"])

// 	// Test case 2: no exclude list, no specified provider, circular queue has no match
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

// }
