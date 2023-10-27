package config

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/paycrest/protocol/utils/logger"
	"github.com/spf13/viper"
)

// OrderConfiguration type defines payment order configurations
type OrderConfiguration struct {
	OrderFulfillmentValidity     time.Duration
	ReceiveAddressValidity       time.Duration
	OrderRequestValidity         time.Duration
	PaycrestOrderContractAddress common.Address
	BundlerRPCURL                string
	PaymasterURL                 string
	EntryPointContractAddress    common.Address
	BucketQueueRebuildInterval   int // in hours
	MaxConcurrentValidators      int
}

// OrderConfig sets the order configuration
func OrderConfig() *OrderConfiguration {
	viper.SetDefault("RECEIVE_ADDRESS_VALIDITY", 30)
	viper.SetDefault("ORDER_REQUEST_VALIDITY", 120)
	viper.SetDefault("ORDER_FULFILLMENT_VALIDITY", 10)
	viper.SetDefault("BUCKET_QUEUE_REBUILD_INTERVAL", 1)
	viper.SetDefault("MAX_CONCURRENT_VALIDATORS", 3)

	return &OrderConfiguration{
		OrderFulfillmentValidity:     time.Duration(viper.GetInt("ORDER_FULFILLMENT_VALIDITY")) * time.Minute,
		ReceiveAddressValidity:       time.Duration(viper.GetInt("RECEIVE_ADDRESS_VALIDITY")) * time.Minute,
		OrderRequestValidity:         time.Duration(viper.GetInt("ORDER_REQUEST_VALIDITY")) * time.Second,
		PaycrestOrderContractAddress: common.HexToAddress(viper.GetString("PAYCREST_ORDER_CONTRACT_ADDRESS")),
		BundlerRPCURL:                viper.GetString("BUNDLER_RPC_URL"),
		PaymasterURL:                 viper.GetString("PAYMASTER_URL"),
		EntryPointContractAddress:    common.HexToAddress(viper.GetString("ENTRY_POINT_CONTRACT_ADDRESS")),
		BucketQueueRebuildInterval:   viper.GetInt("BUCKET_QUEUE_REBUILD_INTERVAL"),
		MaxConcurrentValidators:      viper.GetInt("MAX_CONCURRENT_VALIDATORS"),
	}
}

func init() {
	if err := SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
}
