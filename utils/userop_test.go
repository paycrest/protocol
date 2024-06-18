package utils

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jarcoal/httpmock"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/utils/test"
	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"github.com/stretchr/testify/assert"
)

var (
	PAYMASTER_URL_ETHEREUM = "https://api.stackup.sh/v1/xxx"
	BUNDLER_URL_ETHEREUM   = "https://api.stackup.sh/v1/yyy"
)

func setup() (string, error) {
	// Set up test data
	file, err := test.CreateEnvFile(fmt.Sprintf("%d.env", time.Now().UnixNano()), map[string]string{
		"ACTIVE_AA_SERVICE":      "biconomy",
		"PAYMASTER_URL_ETHEREUM": PAYMASTER_URL_ETHEREUM,
		"BUNDLER_URL_ETHEREUM":   BUNDLER_URL_ETHEREUM,
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
	t.Run("ActiveAAService should change to 'biconomy'", func(t *testing.T) {
		assert.Equal(t, "biconomy", config.OrderConfig().ActiveAAService)
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

	t.Run("test success SendUserOperation", func(t *testing.T) {
		// activate httpmock
		httpmock.Activate()
		defer httpmock.Deactivate()

		data := &userop.UserOperation{
			Sender:               common.Address{}, // Assuming common.Address has a zero value
			Nonce:                big.NewInt(12345),
			InitCode:             []byte{0x01, 0x02, 0x03}, // Example byte slice
			CallData:             []byte{0x04, 0x05, 0x06}, // Example byte slice
			CallGasLimit:         big.NewInt(67890),
			VerificationGasLimit: big.NewInt(111213),
			PreVerificationGas:   big.NewInt(141516),
			MaxFeePerGas:         big.NewInt(171819),
			MaxPriorityFeePerGas: big.NewInt(192021),
			PaymasterAndData:     []byte{0x22, 0x23, 0x24}, // Example byte slice
			Signature:            []byte{0x25, 0x26, 0x27}, // Example byte slice
		}

		// register mock response
		httpmock.RegisterResponder("POST", BUNDLER_URL_ETHEREUM,
			func(r *http.Request) (*http.Response, error) {
				bytes, err := io.ReadAll(r.Body)
				if err != nil {
					log.Fatal(err)
				}

				if strings.Contains(string(bytes), "eth_sendUserOperation") {
					if config.OrderConfig().ActiveAAService == "biconomy" {
						assert.True(t, strings.Contains(string(bytes), "validation_and_execution"))
					} else {
						assert.False(t, strings.Contains(string(bytes), "validation_and_execution"))
					}
					resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      1,
						"result":  "0xa12db69d8d43384e9d9da5f0e9b9698c97c9a90c3447aa815d3f28daf21d3834",
					})
					return resp, err
				} else if strings.Contains(string(bytes), "eth_getUserOperationByHash") {
					resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      1,
						"result": map[string]interface{}{
							"transactionHash": "0xa12db69d8d43384e9d9da5f0e9b9698c97c9a90c3447aa815d3f28daf21d3834",
							"blockNumber":     120,
						},
					})
					return resp, err
				}

				return httpmock.NewBytesResponse(200, []byte(`{"jsonrpc": "2.0","id": 1,"result":[]}`)), nil

			},
		)
		transactionHash, blockNumber, err := SendUserOperation(data, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, transactionHash, "transactionHash should not be empty")
		assert.Equal(t, int64(120), blockNumber)

	})

	t.Run("test success SponsorUserOperation", func(t *testing.T) {
		// activate httpmock
		httpmock.Activate()
		defer httpmock.Deactivate()

		data := &userop.UserOperation{
			Sender:               common.Address{}, // Assuming common.Address has a zero value
			Nonce:                big.NewInt(12345),
			InitCode:             []byte{0x01, 0x02, 0x03}, // Example byte slice
			CallData:             []byte{0x04, 0x05, 0x06}, // Example byte slice
			CallGasLimit:         big.NewInt(67890),
			VerificationGasLimit: big.NewInt(111213),
			PreVerificationGas:   big.NewInt(141516),
			MaxFeePerGas:         big.NewInt(171819),
			MaxPriorityFeePerGas: big.NewInt(192021),
			PaymasterAndData:     []byte{0x22, 0x23, 0x24}, // Example byte slice
			Signature:            []byte{0x25, 0x26, 0x27}, // Example byte slice
		}

		// register mock response
		httpmock.RegisterResponder("POST", PAYMASTER_URL_ETHEREUM,
			func(r *http.Request) (*http.Response, error) {
				bytes, err := io.ReadAll(r.Body)
				if err != nil {
					log.Fatal(err)
				}

				var resp *http.Response
				if strings.Contains(string(bytes), "pm_sponsorUserOperation") {
					if config.OrderConfig().ActiveAAService == "biconomy" {
						assert.True(t, strings.Contains(string(bytes), "INFINITISM"))
						resp, err = httpmock.NewJsonResponse(200, map[string]interface{}{
							"jsonrpc": "2.0",
							"id":      1,
							"result": map[string]interface{}{
								"paymasterAndData":     "0x00000f79b7faf42eebadba19acc07cd08af447890000000000000000000...",
								"preVerificationGas":   "186034",
								"verificationGasLimit": 395693,
								"callGasLimit":         55412,
							},
						})
					} else {
						assert.True(t, strings.Contains(string(bytes), "payg"))
						resp, err = httpmock.NewJsonResponse(200, map[string]interface{}{
							"jsonrpc": "2.0",
							"id":      1,
							"result": map[string]interface{}{
								"paymasterAndData":     "0x00000f79b7faf42eebadba19acc07cd08af447890000000000000000000...",
								"preVerificationGas":   "0x1234",
								"verificationGasLimit": "0x1234",
								"callGasLimit":         "0x1234",
							},
						})
					}
					return resp, err
				}
				return httpmock.NewBytesResponse(200, []byte(`{"jsonrpc": "2.0","id": 1,"result":[]}`)), nil

			},
		)
		err := SponsorUserOperation(data, "sponsored", "", 1)
		assert.NoError(t, err)
	})
}
