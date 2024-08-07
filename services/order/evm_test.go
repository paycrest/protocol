package order

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/enttest"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils/test"
	"github.com/stretchr/testify/assert"
)

var testCxtEVM = struct {
	blockchainClient types.RPCClient
	user             *ent.User
	paymentOrder     *ent.PaymentOrder
	currency         *ent.FiatCurrency
	client           types.RPCClient
	token            *ent.Token
}{}

const (
	rpc = "http://localhost:8545/"
)

func setupEVM() error {

	backend, err := test.SetUpTestBlockchain(map[string]interface{}{
		"networkRPC": rpc,
	})
	if err != nil {
		return err
	}

	testCxtEVM.blockchainClient = backend

	token, err := test.CreateERC20Token(backend, map[string]interface{}{
		"networkRPC":     rpc,
		"identifier":     "localhost_mock",
		"chainID":        int64(1337),
		"deployContract": false,
	})
	if err != nil {
		return fmt.Errorf("evm_test.CreateERC20Token.setup %w", err)
	}

	testCxtEVM.token = token

	user, err := test.CreateTestUser(map[string]interface{}{
		"scope": "provider",
		"email": "providerjohndoe@test.com",
	})
	if err != nil {
		return fmt.Errorf("evm_test.CreateTestUser.setup %w", err)

	}

	testCxtEVM.user = user

	currency, err := test.CreateTestFiatCurrency(map[string]interface{}{
		"code":        "KES",
		"short_name":  "Shilling",
		"decimals":    2,
		"symbol":      "KSh",
		"name":        "Kenyan Shilling",
		"market_rate": 550.0,
	})
	if err != nil {
		return fmt.Errorf("evm_test.CreateTestFiatCurrency.setup %w", err)
	}

	testCxtEVM.currency = currency

	sender, err := test.CreateTestSenderProfile(map[string]interface{}{
		"user_id":     user.ID,
		"currency_id": currency.ID,
	})
	if err != nil {
		return fmt.Errorf("evm_test.CreateTestSenderProfile.setup %w", err)
	}

	paymentOrder, err := test.CreateTestPaymentOrder(backend, token, map[string]interface{}{
		"sender": sender,
	})

	if err != nil {
		return fmt.Errorf("evm_test.CreateTestPaymentOrder.setup %w", err)
	}
	testCxtEVM.paymentOrder = paymentOrder

	return nil

}

func TestEVMOrders(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	err := setupEVM()
	assert.NoError(t, err)

	orderservice := NewOrderEVM()

	t.Run("createOrder", func(t *testing.T) {
		// activate httpmock
		defer httpmock.Deactivate()

		// RPC mock
		httpmock.RegisterResponder("POST", rpc,
			func(r *http.Request) (*http.Response, error) {
				bytes, err := io.ReadAll(r.Body)
				if err != nil {
					log.Fatal(err)
					return httpmock.NewBytesResponse(200, []byte(`{"jsonrpc": "2.0","id": 1,"result":[]}`)), nil

				}

				if strings.Contains(string(bytes), "eth_sendUserOperation") {
					resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      1,
						"result":  "0xa12db69d8d43384e9d9da5f0e9b9698c97c9a90c3447aa815d3f28daf21d3834",
					})
					return resp, fmt.Errorf("evm_test_rpc_mock_eth_sendUserOperation %w", err)
				} else if strings.Contains(string(bytes), "eth_getUserOperationByHash") {
					resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      1,
						"result": map[string]interface{}{
							"transactionHash": "0xa12db69d8d43384e9d9da5f0e9b9698c97c9a90c3447aa815d3f28daf21d3834",
							"blockNumber":     120,
						},
					})
					return resp, fmt.Errorf("eth_getUserOperationByHash %w", err)
				}

				return httpmock.NewBytesResponse(200, []byte(`{"jsonrpc": "2.0","id": 1,"result":[]}`)), nil

			},
		)
		httpmock.RegisterResponder("POST", "http://localhost:8545/rpc",
			func(r *http.Request) (*http.Response, error) {
				bytes, err := io.ReadAll(r.Body)
				if err != nil {
					log.Fatal(err)
					return httpmock.NewBytesResponse(200, []byte(`{"jsonrpc": "2.0","id": 1,"result":[]}`)), nil
				}

				if strings.Contains(string(bytes), "pm_sponsorUserOperation") {
					// if orderConf.ActiveAAService == "biconomy" {
					// 	assert.True(t, strings.Contains(string(bytes), "INFINITISM"))
					// 	resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
					// 		"jsonrpc": "2.0",
					// 		"id":      1,
					// 		"result": map[string]interface{}{
					// 			"paymasterAndData":     "0x00000f79b7faf42eebadba19acc07cd08af447890000000000000000000...",
					// 			"preVerificationGas":   "186034",
					// 			"verificationGasLimit": 395693,
					// 			"callGasLimit":         55412,
					// 		},
					// 	})
					// 	return resp, err
					// } else
					{
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
				} else if strings.Contains(string(bytes), "eth_sendUserOperation") {
					resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      1,
						"result":  "0xa12db69d8d43384e9d9da5f0e9b9698c97c9a90c3447aa815d3f28daf21d3834",
					})
					return resp, fmt.Errorf("evm_test_rpc_mock_eth_sendUserOperation %w", err)
				} else if strings.Contains(string(bytes), "eth_getUserOperationByHash") {
					resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      1,
						"result": map[string]interface{}{
							"transactionHash": "0xa12db69d8d43384e9d9da5f0e9b9698c97c9a90c3447aa815d3f28daf21d3834",
							"blockNumber":     120,
						},
					})
					return resp, fmt.Errorf("eth_getUserOperationByHash %w", err)
				}

				return httpmock.NewBytesResponse(200, []byte(`{"jsonrpc": "2.0","id": 1,"result":[]}`)), nil

			},
		)
		err = orderservice.CreateOrder(context.Background(), testCxtEVM.paymentOrder.ID)
		assert.NoError(t, err)
	})
}
