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
						"logs": []interface{}{
								map[string]interface{}{
									"topics": []interface{}{
										"0x49628fd1471006c1482da88028e9ce4dbb080b815c9b0344d39e5a8e6ec1419f",
										"0x15178e60f1185869b7001608bc4f57ebd9064b77c7f8cc08db1da75464b535b3",
										"0x000000000000000000000000500081e858f7a214cfe4f729f0321ee0e24f900f",
										"0x00000000000000000000000000000f79b7faf42eebadba19acc07cd08af44789",
									},
									"blockNumber": "0x3494a41",
									"transactionHash": "0xd4c6ca8929833176505e4a1ee5693a5f5116415414ab672d0949bbd82df4ff54",
								},
							},
						"receipt": map[string]interface{}{
							"blockNumber":       "0x3494a41",
							"logs": []interface{}{
								map[string]interface{}{
									"address": "0x5ff137d4b0fdcd49dca30c7cf57e578a026d2789",
									"topics": []interface{}{
										"0x40ccd1ceb111a3c186ef9911e1b876dc1f789ed331b86097b3b8851055b6a137",
									},
									"data":             "000000000000000000000000000000000000000000000000000000000000235a71dd0fb221e509dd357a139d634d109422bd2d105c555754afc840ce9f08c8a100000000000000000000000000000000000000000000000000000000000006a30000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000015857716a554b69537231374e4469777449422b486a68365237357371796557577144324b462f4641367338505a4370484e6d5857534e464d6644754b6649676f7a424978414762685862376f6c6d5463367264427575774f654163594c4e4d45564c396465385047582f654c2f796d7a6d50796a6a393535354879625538675949453062337a4c445a64474d314e7a3465557530576e4a2b494946357353424c373678506d5431726b4546797463505a2f666f44534f7368326b7973466879526b634e3356365a717832716957733938484362757a724358614f2b445a634e6836666e492f2b66506f4b35674e74384e684e77714b703830376c354131347755396d6a53625a376978594278324473516f69576677474e616b6559505745445269514d36547a766f33727750442f4d54304c75784b2f472f6556785a7a574c2b74324847336846746748476958745161366f51477363413d3d0000000000000000",
									"blockNumber":      "0x3494a41",
									"transactionHash":  "0xd4c6ca8929833176505e4a1ee5693a5f5116415414ab672d0949bbd82df4ff54",
									"transactionIndex": "0x68",
									"blockHash":        "0x93412643414b664e4fd3467ff58c3bcbd6dc83b84080e1f044e7a16fcf015806",
								},
							},
							"logsBloom":         "0x000000000000010000000000000000000000000000000000000000000000020000090000000000040002001100400000001081000000000000000200000000000000000000000000000000280000008000080000000020000001000000000000000000000a000000000000000000280000000000080000408080001000080000004008000004000000000000000000020000000000000000000000000000000020000000000400000050...",
							"status":            "0x1",
							"to":                "0x5ff137d4b0fdcd49dca30c7cf57e578a026d2789",
							"transactionIndex":  "0x68",
							"type":              "0x2",
						},
					},
				})
				
				return resp, err
				}

				return httpmock.NewBytesResponse(200, []byte(`{"jsonrpc": "2.0","id": 1,"result":[]}`)), nil

			},
		)
		transactionHash, _, blockNumber, err := SendUserOperation(data, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, transactionHash, "transactionHash should not be empty")
		assert.Equal(t, int64(55134785), blockNumber)

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
