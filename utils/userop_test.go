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
									"blockNumber":     "0x708111d", // Hexadecimal string
									"transactionHash": "0xa12db69d8d43384e9d9da5f0e9b9698c97c9a90c3447aa815d3f28daf21d3834",
                                    "address": "0xNotGatewayContractAddress",
                                    "topics": []interface{}{
                                        "0x40ccd1ceb111a3c186ef9911e1b876dc1f789ed331b86097b3b8851055b6a137",
                                    },
                                    "data": "0x0000000000000000000000000000000000000000000000000000000000002710c50a8ae4054535be14a2580a17f860a2177a0077efbda1e531c023e7c7d0efd3000000000000000000000000000000000000000000000000000000000000061400000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000158465852663448713954503853564c76554531674d4a6c6e434933375a342f356c2b76392b32304f77346c556c426f5774593847483164557a776a64694f684f7a32436555316843477664306f6c6550724b574c7a686d41335454367a4870724845732b6e3164456b636f736179676362455a334f4749502f497359364c37575a486a453077536e746d65454a49342f727654447771354b374a4d413461497a5a7163715a78464c75364e4433726b766d5673393935536b42616e4e38767943305650317274356b35467843594d3976734b4a345474626643692b486539384d6b5862794e7435316a5035612f4231664271533066437a587533302f6d634b4449614d776d4571594b703652473351587a73626f7043523055395358462f455273615a7134514c5847392f4e4c476f514747714257334f665033616d496d566b595979506a4a44796863584c5a78667665444b395558673d3d0000000000000000",
                                },
                            },
							"receipt": map[string]interface{}{
								"logs": []interface{}{
								map[string]interface{}{
									"blockNumber":     "0x708111d",
									"transactionHash": "0xa12db69d8d43384e9d9da5f0e9b9698c97c9a90c3447aa815d3f28daf21d3834",
                                    "address": "0xGatewayContractAddress",
                                    "topics": []interface{}{
                                        "0x40ccd1ceb111a3c186ef9911e1b876dc1f789ed331b86097b3b8851055b6a137",
                                    },
                                    "data": "0x0000000000000000000000000000000000000000000000000000000000002710c50a8ae4054535be14a2580a17f860a2177a0077efbda1e531c023e7c7d0efd3000000000000000000000000000000000000000000000000000000000000061400000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000158465852663448713954503853564c76554531674d4a6c6e434933375a342f356c2b76392b32304f77346c556c426f5774593847483164557a776a64694f684f7a32436555316843477664306f6c6550724b574c7a686d41335454367a4870724845732b6e3164456b636f736179676362455a334f4749502f497359364c37575a486a453077536e746d65454a49342f727654447771354b374a4d413461497a5a7163715a78464c75364e4433726b766d5673393935536b42616e4e38767943305650317274356b35467843594d3976734b4a345474626643692b486539384d6b5862794e7435316a5035612f4231664271533066437a587533302f6d634b4449614d776d4571594b703652473351587a73626f7043523055395358462f455273615a7134514c5847392f4e4c476f514747714257334f665033616d496d566b595979506a4a44796863584c5a78667665444b395558673d3d0000000000000000",
                                },
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
