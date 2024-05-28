package utils

import (
	"fmt"
	"os"

	"testing"
	"time"

	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/utils/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setup() (string, error) {
	// Set up test data
	file, err := test.CreateEnvFile(fmt.Sprintf("%d.env", time.Now().UnixNano()), map[string]string{
		"ACTIVE_AA_SERVICE":      "BICONOMY",
		"PAYMASTER_URL_ETHEREUM": "https://api.stackup.sh/v1/xxx",
		"BUNDLER_URL_ETHEREUM":   "https://api.stackup.sh/v1/yyy",
	})
	if err != nil {
		return "", err
	}

	// change .env file loaded
	viper.AddConfigPath(".")
	viper.SetConfigName(file)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	return file, nil
}
func TestUserop(t *testing.T) {
	file, err := setup()
	assert.NoError(t, err)

	// remove the generated file
	t.Cleanup(func() {
		os.Remove(file)
	})
	// TEST TODO
	t.Run("ActiveAAService should change to BICONOMY", func(t *testing.T) {
		assert.Equal(t, "BICONOMY", config.OrderConfig().ActiveAAService)
	})

	t.Run("test getEndpoints", func(t *testing.T) {
		t.Run("when chainID is supported getEndpoints", func(t *testing.T) {
			bundlerID, paymaster, err := getEndpoints(1)
			assert.NoError(t, err)
			assert.NotEmpty(t, bundlerID, "bundlerID should not be empty")
			assert.NotEmpty(t, paymaster, "paymaster should not be empty")
		})

		t.Run("when chainID is not supported getEndpoints", func(t *testing.T) {
			bundlerID, paymaster, err := getEndpoints(1000)
			assert.Error(t, err)
			assert.Empty(t, bundlerID, "bundlerID should be empty")
			assert.Empty(t, paymaster, "paymaster should be empty")
		})

	})
}
