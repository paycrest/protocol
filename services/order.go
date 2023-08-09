package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/paycrest/paycrest-protocol/config"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/types"
	"github.com/paycrest/paycrest-protocol/utils"
	cryptoUtils "github.com/paycrest/paycrest-protocol/utils/crypto"
)

// OrderService provides functionality related to on-chain interactions for payment orders
type OrderService struct{}

type CreateOrderParams struct {
	Token              common.Address
	Amount             *big.Int
	RefundAddress      common.Address
	SenderFeeRecipient common.Address
	SenderFee          *big.Int
	Rate               *big.Int
	InstitutionCode    [32]byte
	MessageHash        string
}

// UserOperation represents an EIP-4337 style transaction for a smart contract account.
type UserOperation struct {
	Sender               common.Address `json:"sender"               mapstructure:"sender"               validate:"required"`
	Nonce                *big.Int       `json:"nonce"                mapstructure:"nonce"                validate:"required"`
	InitCode             []byte         `json:"initCode"             mapstructure:"initCode"             validate:"required"`
	CallData             []byte         `json:"callData"             mapstructure:"callData"             validate:"required"`
	CallGasLimit         *big.Int       `json:"callGasLimit"         mapstructure:"callGasLimit"         validate:"required"`
	VerificationGasLimit *big.Int       `json:"verificationGasLimit" mapstructure:"verificationGasLimit" validate:"required"`
	PreVerificationGas   *big.Int       `json:"preVerificationGas"   mapstructure:"preVerificationGas"   validate:"required"`
	MaxFeePerGas         *big.Int       `json:"maxFeePerGas"         mapstructure:"maxFeePerGas"         validate:"required"`
	MaxPriorityFeePerGas *big.Int       `json:"maxPriorityFeePerGas" mapstructure:"maxPriorityFeePerGas" validate:"required"`
	PaymasterAndData     []byte         `json:"paymasterAndData"     mapstructure:"paymasterAndData"     validate:"required"`
	Signature            []byte         `json:"signature"            mapstructure:"signature"            validate:"required"`
}

// NewOrderService creates a new instance of OrderService.
func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, client types.RPCClient, order *ent.PaymentOrder) error {
	var conf = config.OrderConfig()
	fromAddress, _, _ := utils.GetMasterAccount()

	// Connect to RPC endpoint
	var err error
	if client == nil {
		client, err = types.NewEthClient(order.Edges.Token.Edges.Network.RPCEndpoint)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	// Create createOrder data
	data, err := s.createOrderCallData(order)
	if err != nil {
		return fmt.Errorf("failed to create createOrder calldata: %w", err)
	}

	call := ethereum.CallMsg{
		To:    &conf.PaycrestOrderContractAddress,
		Value: big.NewInt(0),
		Data:  data,
	}

	executeABI, err := abi.JSON(strings.NewReader("execute(address,uint256,bytes)"))
	if err != nil {
		return fmt.Errorf("failed to create execute ABI: %w", err)
	}

	calldata, err := executeABI.Pack("execute", conf.PaycrestOrderContractAddress, big.NewInt(0), data)
	if err != nil {
		return fmt.Errorf("failed to pack execute ABI: %w", err)
	}

	callGasLimit, err := client.EstimateGas(ctx, call)
	if err != nil {
		return fmt.Errorf("failed to estimate gas: %w", err)
	}

	nonce, err := client.PendingNonceAt(ctx, *fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	// TODO: build user operation
	userOperation := &UserOperation{
		Sender:               common.HexToAddress(order.Edges.ReceiveAddress.Address),
		Nonce:                big.NewInt(int64(nonce)),
		InitCode:             common.FromHex("0x"),
		CallData:             calldata,
		CallGasLimit:         big.NewInt(int64(callGasLimit)),
		VerificationGasLimit: big.NewInt(70000),
		PreVerificationGas:   big.NewInt(21000),
		MaxFeePerGas:         gasPrice,
		MaxPriorityFeePerGas: gasPrice,
		PaymasterAndData:     common.FromHex("0x"),
	}

	_ = userOperation

	// TODO: send user operation

	return nil
}

// createOrderCallData creates the data for the createOrder method
func (s *OrderService) createOrderCallData(order *ent.PaymentOrder) ([]byte, error) {
	fromAddress, _, _ := utils.GetMasterAccount()

	// Encrypt recipient details
	encryptedOrderRecipient, err := s.encryptOrderRecipient(order.Edges.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt recipient details: %w", err)
	}

	// Define params
	params := &CreateOrderParams{
		Token:              common.HexToAddress(order.Edges.Token.ContractAddress),
		Amount:             order.Amount.BigInt(),
		RefundAddress:      *fromAddress,
		SenderFeeRecipient: *fromAddress,
		SenderFee:          big.NewInt(0),
		Rate:               big.NewInt(0),
		InstitutionCode:    utils.StringTo32Byte(order.Edges.Recipient.Institution),
		MessageHash:        encryptedOrderRecipient,
	}

	createOrderABI, err := abi.JSON(strings.NewReader("createOrder(address,uint256,address,address,uint256,uint96,bytes32,string)"))
	if err != nil {
		return nil, fmt.Errorf("failed to create createOrder ABI: %w", err)
	}

	// Generate call data
	data, err := createOrderABI.Pack("createOrder", params.Token, params.Amount, params.RefundAddress, params.SenderFeeRecipient, params.SenderFee, params.Rate, params.InstitutionCode, params.MessageHash)

	return data, err
}

func (s *OrderService) encryptOrderRecipient(recipient *ent.PaymentOrderRecipient) (string, error) {
	message := []struct {
		AccountIdentifier string
		AccountName       string
		Institution       string
	}{
		{recipient.AccountIdentifier, recipient.AccountName, recipient.Institution},
	}
	messageCipher, err := cryptoUtils.EncryptJSON(message)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt message: %w", err)
	}

	return fmt.Sprintf("0x%x", messageCipher), nil
}
