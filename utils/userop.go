package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/services/contracts"
	"github.com/paycrest/protocol/types"
	cryptoUtils "github.com/paycrest/protocol/utils/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	fromAddress, privateKey, _ = cryptoUtils.GenerateAccountFromIndex(0)
	OrderConf                  = config.OrderConfig()
)

// Initialize user operation with defaults
func InitializeUserOperation(ctx context.Context, client types.RPCClient, rpcUrl, sender, salt string) (*userop.UserOperation, error) {
	var err error

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(rpcUrl)
		if err != nil {
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
	nonce, err := client.PendingNonceAt(ctx, *fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}
	userOperation.Nonce = big.NewInt(int64(nonce))

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
	maxFeePerGas, maxPriorityFeePerGas := eip1559GasPrice(ctx, client)
	userOperation.MaxFeePerGas = maxFeePerGas
	userOperation.MaxPriorityFeePerGas = maxPriorityFeePerGas

	return userOperation, nil
}

// SponsorUserOperation sponsors the user operation
// ref: https://docs.stackup.sh/docs/paymaster-api-rpc-methods#pm_sponsoruseroperation
func SponsorUserOperation(userOp *userop.UserOperation, mode string) error {
	client, err := rpc.Dial(OrderConf.PaymasterURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	var payload map[string]interface{}

	switch mode {
	case "payg":
		payload = map[string]interface{}{
			"type": "payg",
		}
	case "erc20token":
		payload = map[string]interface{}{
			"type":  "erc20token",
			"token": "0x9999f7Fea5938fD3b1E26A12c3f2fb024e194f97",
		}
	default:
		return fmt.Errorf("invalid mode")
	}

	requestParams := []interface{}{
		userOp,
		OrderConf.EntryPointContractAddress.Hex(),
		payload,
	}

	// op, _ := userOp.MarshalJSON()
	// fmt.Println(string(op))

	var result json.RawMessage
	err = client.Call(&result, "pm_sponsorUserOperation", requestParams...)
	if err != nil {
		return fmt.Errorf("RPC error: %w", err)
	}

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

	return nil
}

// SignUserOperation signs the user operation
func SignUserOperation(userOperation *userop.UserOperation) error {
	// Sign user operation
	userOpHash := userOperation.GetUserOpHash(
		OrderConf.EntryPointContractAddress,
		big.NewInt(137),
	)

	signature, err := PersonalSign(string(userOpHash[:]), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign user operation: %w", err)
	}
	userOperation.Signature = signature

	return nil
}

// SendUserOperation sends the user operation
func SendUserOperation(userOp *userop.UserOperation) (string, error) {
	client, err := rpc.Dial(OrderConf.BundlerRPCURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	requestParams := []interface{}{
		userOp,
		OrderConf.EntryPointContractAddress.Hex(),
	}

	// op, _ := userOp.MarshalJSON()
	// fmt.Println(string(op))

	var result json.RawMessage
	err = client.Call(&result, "eth_sendUserOperation", requestParams...)
	if err != nil {
		return "", fmt.Errorf("RPC error: %w", err)
	}

	var userOpHash string
	err = json.Unmarshal(result, &userOpHash)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return userOpHash, nil
}

// GetUserOperationStatus returns the status of the user operation
func GetUserOperationStatus(userOpHash string) (bool, error) {
	client, err := rpc.Dial(OrderConf.BundlerRPCURL)
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
func eip1559GasPrice(ctx context.Context, client types.RPCClient) (maxFeePerGas, maxPriorityFeePerGas *big.Int) {
	tip, _ := client.SuggestGasTipCap(ctx)
	latestHeader, _ := client.HeaderByNumber(ctx, nil)

	buffer := new(big.Int).Mul(tip, big.NewInt(13)).Div(tip, big.NewInt(100))
	maxPriorityFeePerGas = new(big.Int).Add(tip, buffer)

	if latestHeader.BaseFee != nil {
		maxFeePerGas = new(big.Int).
			Mul(latestHeader.BaseFee, big.NewInt(2)).
			Add(latestHeader.BaseFee, maxPriorityFeePerGas)
	} else {
		maxFeePerGas = maxPriorityFeePerGas
	}

	return maxFeePerGas, maxPriorityFeePerGas
}
