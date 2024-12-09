package utils

import (
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"

	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jarcoal/httpmock"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"github.com/stretchr/testify/assert"
)

func TestUserOp(t *testing.T) {

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
		httpmock.RegisterResponder("POST", orderConf.BundlerUrlEthereum,
			func(r *http.Request) (*http.Response, error) {
				bytes, err := io.ReadAll(r.Body)
				if err != nil {
					log.Fatal(err)
				}

				if strings.Contains(string(bytes), "eth_sendUserOperation") {
					if orderConf.ActiveAAService == "biconomy" {
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
				} else if strings.Contains(string(bytes), "eth_getUserOperationReceipt") {
					resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
                        "jsonrpc": "2.0",
                        "id":      1,
                        "result": map[string]interface{}{
                            "transactionHash": "0xa12db69d8d43384e9d9da5f0e9b9698c97c9a90c3447aa815d3f28daf21d3834",
                            "blockNumber":     "0x708111d", // Hexadecimal string
                            "logs": []interface{}{
                                map[string]interface{}{
                                    "address": "0xYourGatewayContractAddress",
                                    "topics": []interface{}{
                                        "0x40ccd1ceb111a3c186ef9911e1b876dc1f789ed331b86097b3b8851055b6a137",
                                    },
                                    "data": "0x000000000000000000000000000000000000000000000000000000000000206f89fb9887446bc159462a4bb070d9b809f0e3474b0df07b1b664e1f19a74513fd000000000000000000000000000000000000000000000000000000000000069000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000158545041777a45466a673554774774354a44314e5a696d6d6f384b78456c683534633965417766725a6e3946534f67313341444d55655a745a43396b4a4d50344b6544636553556549564d4d576d46376d494e6d382f364674365244786a334e4836423566536b3578544e6775654f6e6f35505578307a39716b584d356b43574f475a544a4c314139617557594b6f4a6b4f59504c49513343384f6f655a6e3053644f7930616679484c7a52707a5478576131615447734b6a3343553358544778747473526c75596931456439622b73477470466b486d6d3172384f37307867776b6330757648324564703078457a6662495a3657767239755755486841674b5a65416b4259387838686b544d75575470534a693475446670705247414870446850624e51517974777845567852312f445750467137483558547946713879615171454a7867666672594d526254744c554554646333773d3d0000000000000000",
                                },
                            },
                        },
                    })
					return resp, err
				}

				return httpmock.NewBytesResponse(200, []byte(`{"jsonrpc": "2.0","id": 1,"result":[]}`)), nil

			},
		)
		transactionHash, orderId, blockNumber, err := SendUserOperation(data, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, transactionHash, "transactionHash should not be empty")
		assert.NotEmpty(t, orderId, "orderId should not be empty")
		assert.Equal(t, int64(118135645), blockNumber)

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
		httpmock.RegisterResponder("POST", orderConf.PaymasterUrlEthereum,
			func(r *http.Request) (*http.Response, error) {
				bytes, err := io.ReadAll(r.Body)
				if err != nil {
					log.Fatal(err)
				}

				if strings.Contains(string(bytes), "pm_sponsorUserOperation") {
					if orderConf.ActiveAAService == "biconomy" {
						assert.True(t, strings.Contains(string(bytes), "INFINITISM"))
						resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
							"jsonrpc": "2.0",
							"id":      1,
							"result": map[string]interface{}{
								"paymasterAndData":     "0x00000f79b7faf42eebadba19acc07cd08af447890000000000000000000...",
								"preVerificationGas":   "186034",
								"verificationGasLimit": 395693,
								"callGasLimit":         55412,
							},
						})
						return resp, err
					} else {
						assert.True(t, strings.Contains(string(bytes), "payg"))
						resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
							"jsonrpc": "2.0",
							"id":      1,
							"result": map[string]interface{}{
								"paymasterAndData":     "0x00000f79b7faf42eebadba19acc07cd08af447890000000000000000000...",
								"preVerificationGas":   "0x1234",
								"verificationGasLimit": "0x1234",
								"callGasLimit":         "0x1234",
							},
						})
						return resp, err
					}
				}
				return httpmock.NewBytesResponse(200, []byte(`{"jsonrpc": "2.0","id": 1,"result":[]}`)), nil

			},
		)

		err := SponsorUserOperation(data, "sponsored", "", 1)
		assert.NoError(t, err)
	})
}
