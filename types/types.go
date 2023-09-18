package types

import (
	"context"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/shopspring/decimal"
)

// RPCClient is an interface for interacting with the blockchain.
type RPCClient interface {
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	Commit() common.Hash
}

// Custom type that implements RPCClient
type ethRPC struct {
	*ethclient.Client
}

// Implements Commit() method
func (e *ethRPC) Commit() common.Hash {
	return common.Hash{} // no-op
}

// Helper function to create client
func NewEthClient(endpoint string) (RPCClient, error) {

	ethClient, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	return &ethRPC{ethClient}, nil
}

// RegisterPayload is the payload for the register endpoint
type RegisterPayload struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6,max=20"`
	TradingName string `json:"tradingName"`
	Country     string `json:"country"`
}

// RegisterResponse is the response for the register endpoint
type RegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

// LoginPayload is the payload for the login endpoint
type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// LoginResponse is the response for the login endpoint
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshJWTPayload is the payload for the refresh endpoint
type RefreshJWTPayload struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshResponse is the response for the refresh endpoint
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

// CreateAPIKeyPayload is the payload for the generate API key endpoint
type CreateAPIKeyPayload struct {
	Name  string       `json:"name" binding:"required"`
	Scope apikey.Scope `json:"scope" binding:"required,oneof=sender provider tx_validator"`
}

// APIKeyResponse is the response type for an API key
type APIKeyResponse struct {
	ID        uuid.UUID    `json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	Name      string       `json:"name"`
	Scope     apikey.Scope `json:"scope"`
	Secret    string       `json:"secret"`
	IsActive  bool         `json:"isActive"`
}

// ERC20Transfer is the Transfer event of an ERC20 smart contract
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// PaymentOrderRecipient describes a payment order recipient
type PaymentOrderRecipient struct {
	Institution       string `json:"institution" binding:"required"`
	AccountIdentifier string `json:"accountIdentifier" binding:"required"`
	AccountName       string `json:"accountName" binding:"required"`
	ProviderID        string `json:"providerId"`
}

// NewPaymentOrderPayload is the payload for the create payment order endpoint
type NewPaymentOrderPayload struct {
	Amount    decimal.Decimal       `json:"amount" binding:"required"`
	Token     string                `json:"token" binding:"required"`
	Network   string                `json:"network" binding:"required"`
	Recipient PaymentOrderRecipient `json:"recipient" binding:"required"`
}

// ReceiveAddressResponse is the response type for a receive address
type ReceiveAddressResponse struct {
	ID             uuid.UUID `json:"id"`
	Amount         float64   `json:"amount"`
	Network        string    `json:"network"`
	ReceiveAddress string    `json:"receiveAddress"`
}

type PaymentOrderResponse struct {
	ID        uuid.UUID             `json:"id"`
	Amount    float64               `json:"amount"`
	Network   string                `json:"network"`
	Recipient PaymentOrderRecipient `json:"recipient"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	TxHash    string                `json:"txHash"`
	Status    string                `json:"status"`
}

// ConfirmEmailPayload is the payload for the confirmEmail endpoint
type ConfirmEmailPayload struct {
	Token string `json:"token" binding:"required"`
}

// SendEmailPayload is content of a email request.
type SendEmailPayload struct {
	FromAddress string
	ToAddress   string
	Subject     string
	Body        string
}

// SendEmailResponse is the mailgunv3.Send response struct
type SendEmailResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}
