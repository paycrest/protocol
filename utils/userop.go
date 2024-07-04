package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shopspring/decimal"

	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/services/contracts"
	"github.com/paycrest/protocol/types"
	cryptoUtils "github.com/paycrest/protocol/utils/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	fromAddress, privateKey, _ = cryptoUtils.GenerateAccountFromIndex(0)
	orderConf                  = config.OrderConfig()
)

// Initialize user operation with defaults
func InitializeUserOperation(ctx context.Context, client types.RPCClient, rpcUrl, sender, salt string) (*userop.UserOperation, error) {
	var err error

	// Connect to RPC endpoint
	if client == nil {
		retryErr := Retry(3, 5*time.Second, func() error {
			client, err = types.NewEthClient(rpcUrl)
			return err
		})
		if retryErr != nil {
			return nil, fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	// Build user operation
	userOperation := &userop.UserOperation{
		Sender:               common.HexToAddress(sender),
		Nonce:                big.NewInt(0),
		InitCode:             common.FromHex("0x"),
		CallData:             common.FromHex("0x"),
		CallGasLimit:         big.NewInt(350000),
		VerificationGasLimit: big.NewInt(300000),
		PreVerificationGas:   big.NewInt(100000),
		MaxFeePerGas:         big.NewInt(50000),
		MaxPriorityFeePerGas: big.NewInt(1000),
		PaymasterAndData:     common.FromHex("0x"),
		Signature:            common.FromHex("0xa925dcc5e5131636e244d4405334c25f034ebdd85c0cb12e8cdb13c15249c2d466d0bade18e2cafd3513497f7f968dcbb63e519acd9b76dcae7acd61f11aa8421b"),
	}

	// Get nonce
	nonce, err := getNonce(client, userOperation.Sender)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}
	userOperation.Nonce = nonce

	// Create initcode
	code, err := client.CodeAt(ctx, userOperation.Sender, nil)
	if err != nil {
		return nil, err
	}

	if len(code) == 0 {
		// address does not exist yet
		salt, _ := new(big.Int).SetString(salt, 10)

		createAccountCallData, err := createAccountCallData(*fromAddress, salt)
		if err != nil {
			return nil, fmt.Errorf("failed to create init code: %w", err)
		}

		var factoryAddress [20]byte
		copy(factoryAddress[:], common.HexToAddress("0x9406Cc6185a346906296840746125a0E44976454").Bytes())

		userOperation.InitCode = append(factoryAddress[:], createAccountCallData...)
	}

	// Set gas fees
	maxFeePerGas, maxPriorityFeePerGas, err := eip1559GasPrice(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	userOperation.MaxFeePerGas = maxFeePerGas
	userOperation.MaxPriorityFeePerGas = maxPriorityFeePerGas

	return userOperation, nil
}

// SponsorUserOperation sponsors the user operation from stackup
// ref: https://docs.stackup.sh/docs/paymaster-api-rpc-methods#pm_sponsoruseroperation
func SponsorUserOperation(userOp *userop.UserOperation, mode string, token string, chainId int64) error {
	_, paymasterUrl, err := getEndpoints(chainId)
	if err != nil {
		return fmt.Errorf("failed to get endpoints: %w", err)
	}

	client, err := rpc.Dial(paymasterUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	var payload map[string]interface{}
	var requestParams []interface{}

	if orderConf.ActiveAAService == "stackup" {
		switch mode {
		case "sponsored":
			payload = map[string]interface{}{
				"type": "payg",
			}
		case "erc20":
			if token == "" {
				return fmt.Errorf("token address is required")
			}

			payload = map[string]interface{}{
				"type":  "erc20token",
				"token": token,
			}
		default:
			return fmt.Errorf("invalid mode")
		}

		requestParams = []interface{}{
			userOp,
			orderConf.EntryPointContractAddress.Hex(),
			payload,
		}
	} else if orderConf.ActiveAAService == "biconomy" {
		mode = "sponsored"

		switch mode {
		case "sponsored":
			payload = map[string]interface{}{
				"mode": "SPONSORED",
				"sponsorshipInfo": map[string]interface{}{
					"webhookData": map[string]interface{}{},
					"smartAccountInfo": map[string]string{
						"name":    "INFINITISM",
						"version": "1.0.0",
					},
				},
				"expiryDuration":     300,
				"calculateGasLimits": true,
			}
		case "erc20":
			if token == "" {
				return fmt.Errorf("token address is required")
			}

			payload = map[string]interface{}{
				"mode": "ERC20",
				"tokenInfo": map[string]string{
					"feeTokenAddress": token,
				},
				"calculateGasLimits": true,
			}
		default:
			return fmt.Errorf("invalid mode")
		}

		requestParams = []interface{}{
			map[string]interface{}{
				"sender":               userOp.Sender.Hex(),
				"nonce":                userOp.Nonce.String(),
				"initCode":             hexutil.Encode(userOp.InitCode),
				"callData":             hexutil.Encode(userOp.CallData),
				"callGasLimit":         userOp.CallGasLimit.String(),
				"verificationGasLimit": userOp.VerificationGasLimit.String(),
				"preVerificationGas":   userOp.PreVerificationGas.String(),
				"maxFeePerGas":         userOp.MaxFeePerGas.String(),
				"maxPriorityFeePerGas": userOp.MaxPriorityFeePerGas.String(),
				"paymasterAndData":     hexutil.Encode(userOp.PaymasterAndData),
				"signature":            hexutil.Encode(userOp.Signature),
			},
			payload,
		}
	}

	var result json.RawMessage
	err = client.Call(&result, "pm_sponsorUserOperation", requestParams...)
	if err != nil {
		op, _ := userOp.MarshalJSON()
		return fmt.Errorf("RPC error: %w\nUser Operation: %s", err, string(op))
	}

	if orderConf.ActiveAAService == "stackup" {
		type Response struct {
			PaymasterAndData     string `json:"paymasterAndData"     mapstructure:"paymasterAndData"`
			PreVerificationGas   string `json:"preVerificationGas"   mapstructure:"preVerificationGas"`
			VerificationGasLimit string `json:"verificationGasLimit" mapstructure:"verificationGasLimit"`
			CallGasLimit         string `json:"callGasLimit"         mapstructure:"callGasLimit"`
		}
		var response Response

		err = json.Unmarshal(result, &response)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		userOp.CallGasLimit, _ = new(big.Int).SetString(response.CallGasLimit, 0)
		userOp.VerificationGasLimit, _ = new(big.Int).SetString(response.VerificationGasLimit, 0)
		userOp.PreVerificationGas, _ = new(big.Int).SetString(response.PreVerificationGas, 0)
		userOp.PaymasterAndData = common.FromHex(response.PaymasterAndData)

	} else if orderConf.ActiveAAService == "biconomy" {
		var response map[string]interface{}

		err = json.Unmarshal(result, &response)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		userOp.PaymasterAndData = common.FromHex(response["paymasterAndData"].(string))
		userOp.PreVerificationGas, _ = new(big.Int).SetString(response["preVerificationGas"].(string), 0)
		userOp.VerificationGasLimit = decimal.NewFromFloat(response["verificationGasLimit"].(float64)).BigInt()
		userOp.CallGasLimit = decimal.NewFromFloat(response["callGasLimit"].(float64)).BigInt()
	}

	return nil
}

// SignUserOperation signs the user operation
func SignUserOperation(userOperation *userop.UserOperation, chainId int64) error {
	// Sign user operation
	userOpHash := userOperation.GetUserOpHash(
		orderConf.EntryPointContractAddress,
		big.NewInt(chainId),
	)

	signature, err := PersonalSign(string(userOpHash[:]), privateKey)
	if err != nil {
		return err
	}
	userOperation.Signature = signature

	return nil
}

// SendUserOperation sends the user operation
func SendUserOperation(userOp *userop.UserOperation, chainId int64) (string, int64, error) {
	bundlerUrl, _, err := getEndpoints(chainId)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get endpoints: %w", err)
	}

	client, err := rpc.Dial(bundlerUrl)
	if err != nil {
		return "", 0, fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	var requestParams []interface{}

	if orderConf.ActiveAAService == "stackup" {
		requestParams = []interface{}{
			userOp,
			orderConf.EntryPointContractAddress.Hex(),
		}
	} else if orderConf.ActiveAAService == "biconomy" {
		requestParams = []interface{}{
			userOp,
			orderConf.EntryPointContractAddress.Hex(),
			map[string]string{
				"simulation_type": "validation_and_execution",
			},
		}
	}

	var result json.RawMessage
	err = client.Call(&result, "eth_sendUserOperation", requestParams...)
	if err != nil {
		op, _ := userOp.MarshalJSON()
		return "", 0, fmt.Errorf("RPC error: %w\nUser Operation: %s", err, string(op))
	}

	var userOpHash string
	err = json.Unmarshal(result, &userOpHash)
	if err != nil {
		return "", 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	response, err := GetUserOperationByHash(userOpHash, chainId)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get user operation by hash: %w", err)
	}

	transactionHash, ok := response["transactionHash"].(string)
	if !ok {
		return "", 0, fmt.Errorf("failed to get transaction hash")
	}

	blockNumber, ok := response["blockNumber"].(float64)
	if !ok {
		return "", 0, fmt.Errorf("failed to get block number")
	}

	return transactionHash, int64(blockNumber), nil
}

// GetUserOperationByHash fetches the user operation by hash
func GetUserOperationByHash(userOpHash string, chainId int64) (map[string]interface{}, error) {
	bundlerUrl, _, err := getEndpoints(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %w", err)
	}

	client, err := rpc.Dial(bundlerUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	start := time.Now()
	timeout := 10 * time.Minute

	var response map[string]interface{}
	for {
		time.Sleep(10 * time.Second)
		var result json.RawMessage
		err = client.Call(&result, "eth_getUserOperationByHash", []interface{}{userOpHash}...)
		if err != nil {
			return nil, fmt.Errorf("RPC error: %w", err)
		}

		err = json.Unmarshal(result, &response)
		if err != nil {
			return nil, err
		}

		if response == nil && response["transactionHash"] == nil {
			elapsed := time.Since(start)
			if elapsed >= timeout {
				return nil, err
			}
			continue
		}

		break
	}

	return response, nil
}

// GetPaymasterAccount fetches the paymaster account from stackup
// ref: https://docs.stackup.sh/docs/paymaster-api-rpc-methods#pm_accounts
func GetPaymasterAccount(chainId int64) (string, error) {
	if orderConf.ActiveAAService == "biconomy" {
		return "0x00000f79b7faf42eebadba19acc07cd08af44789", nil
	}

	_, paymasterUrl, err := getEndpoints(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get endpoints: %w", err)
	}

	client, err := rpc.Dial(paymasterUrl)
	if err != nil {
		return "", fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	requestParams := []interface{}{
		orderConf.EntryPointContractAddress.Hex(),
	}

	var result json.RawMessage
	err = client.Call(&result, "pm_accounts", requestParams...)
	if err != nil {
		return "", fmt.Errorf("RPC error: %w", err)
	}

	var response []string
	err = json.Unmarshal(result, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response[0], nil
}

// GetUserOperationStatus returns the status of the user operation
func GetUserOperationStatus(userOpHash string, chainId int64) (bool, error) {
	bundlerUrl, _, err := getEndpoints(chainId)
	if err != nil {
		return false, fmt.Errorf("failed to get endpoints: %w", err)
	}

	client, err := rpc.Dial(bundlerUrl)
	if err != nil {
		return false, fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	requestParams := []interface{}{
		userOpHash,
	}

	var result json.RawMessage
	err = client.Call(&result, "eth_getUserOperationReceipt", requestParams)
	if err != nil {
		return false, fmt.Errorf("RPC error: %w", err)
	}

	var userOpStatus map[string]interface{}
	err = json.Unmarshal(result, &userOpStatus)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return userOpStatus["success"].(bool), nil
}

// createAccountCallData creates the data for the createAccount method
func createAccountCallData(owner common.Address, salt *big.Int) ([]byte, error) {
	// Create ABI
	accountFactoryABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountFactoryMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse account factory ABI: %w", err)
	}

	// Create calldata
	calldata, err := accountFactoryABI.Pack("createAccount", owner, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to pack createAccount ABI: %w", err)
	}

	return calldata, nil
}

// eip1559GasPrice computes the EIP1559 gas price
func eip1559GasPrice(ctx context.Context, client types.RPCClient) (maxFeePerGas, maxPriorityFeePerGas *big.Int, err error) {
	latestHeader, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	if latestHeader.BaseFee != nil {
		tip, err := client.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, nil, err
		}
		maxFeePerGas = big.NewInt(0).Add(tip, new(big.Int).Mul(latestHeader.BaseFee, common.Big2))
		maxPriorityFeePerGas = tip
	} else {
		sgp, err := client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, nil, err
		}
		maxFeePerGas = sgp
		maxPriorityFeePerGas = sgp
	}

	return maxFeePerGas, maxPriorityFeePerGas, nil
}

// getEndpoints returns the bundler and paymaster URLs for the given chain ID
func getEndpoints(chainId int64) (bundlerUrl, paymasterUrl string, err error) {
	switch chainId {
	case 1:
		bundlerUrl = orderConf.BundlerUrlEthereum
		paymasterUrl = orderConf.PaymasterUrlEthereum
	case 11155111:
		bundlerUrl = orderConf.BundlerUrlEthereum
		paymasterUrl = orderConf.PaymasterUrlEthereum
	case 137:
		bundlerUrl = orderConf.BundlerUrlPolygon
		paymasterUrl = orderConf.PaymasterUrlPolygon
	case 56:
		bundlerUrl = orderConf.BundlerUrlBSC
		paymasterUrl = orderConf.PaymasterUrlBSC
	case 8453:
		bundlerUrl = orderConf.BundlerUrlBase
		paymasterUrl = orderConf.PaymasterUrlBase
	case 84532:
		bundlerUrl = orderConf.BundlerUrlBase
		paymasterUrl = orderConf.PaymasterUrlBase
	case 42161:
		bundlerUrl = orderConf.BundlerUrlArbitrum
		paymasterUrl = orderConf.PaymasterUrlArbitrum
	case 421614:
		bundlerUrl = orderConf.BundlerUrlArbitrum
		paymasterUrl = orderConf.PaymasterUrlArbitrum
	default:
		return "", "", fmt.Errorf("unsupported chain ID")
	}

	return bundlerUrl, paymasterUrl, nil
}

// getNonce returns the nonce for the given sender
// https://docs.stackup.sh/docs/useroperation-nonce
func getNonce(client types.RPCClient, sender common.Address) (nonce *big.Int, err error) {
	entrypoint, err := contracts.NewEntryPoint(orderConf.EntryPointContractAddress, client.(bind.ContractBackend))
	if err != nil {
		return nil, err
	}

	key := big.NewInt(0)
	nonce, err = entrypoint.GetNonce(nil, sender, key)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}
