package config

import (
	"time"

	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/spf13/viper"
)

// OrderConfiguration type defines payment order configurations
type OrderConfiguration struct {
	ReceiveAddressValidity time.Duration
}

// OrderConfig sets the order configuration
func OrderConfig() *OrderConfiguration {
	viper.SetDefault("RECEIVE_ADDRESS_VALIDITY", 30)

	return &OrderConfiguration{
		ReceiveAddressValidity: time.Duration(viper.GetInt("RECEIVE_ADDRESS_VALIDITY")) * time.Minute,
	}
}

func init() {
	if err := SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
}
