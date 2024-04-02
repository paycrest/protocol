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
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/lockorderfulfillment"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/providerordertoken"
	"github.com/paycrest/protocol/ent/providerprofile"
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
	FirstName   string   `json:"firstName" binding:"required"`
	LastName    string   `json:"lastName" binding:"required"`
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6,max=20"`
	TradingName string   `json:"tradingName"`
	Currency    string   `json:"currency"`
	Scopes      []string `json:"scopes" binding:"required,dive,oneof=sender provider"`
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

// LockOrderResponse is the response for the lock payment order model
type LockOrderResponse struct {
	ID                uuid.UUID               `json:"id"`
	Amount            decimal.Decimal         `json:"amount"`
	Token             string                  `json:"token"`
	Institution       string                  `json:"institution"`
	AccountIdentifier string                  `json:"accountIdentifier"`
	AccountName       string                  `json:"accountName"`
	Status            lockpaymentorder.Status `json:"status"`
	UpdatedAt         time.Time               `json:"updatedAt"`
}

// AcceptOrderResponse is the response for the accept order endpoint
type AcceptOrderResponse struct {
	ID                uuid.UUID       `json:"id"`
	Amount            decimal.Decimal `json:"amount"`
	Institution       string          `json:"institution"`
	AccountIdentifier string          `json:"accountIdentifier"`
	AccountName       string          `json:"accountName"`
	Memo              string          `json:"memo"`
}

// FulfillLockOrderPayload is the payload for the fulfill order endpoint
type FulfillLockOrderPayload struct {
	TxID             string                                `json:"txId" binding:"required"`
	ValidationStatus lockorderfulfillment.ValidationStatus `json:"validationStatus"`
	ValidationError  string                                `json:"validationError"`
}

// CancelLockOrderPayload is the payload for the cancel order endpoint
type CancelLockOrderPayload struct {
	Reason string `json:"reason" binding:"required"`
}

// LoginPayload is the payload for the login endpoint
type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// LoginResponse is the response for the login endpoint
type LoginResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	Scopes       []string `json:"scopes"`
}

// RefreshJWTPayload is the payload for the refresh endpoint
type RefreshJWTPayload struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// ValidatorProfilePayload is the payload for the validator profile endpoint
type ValidatorProfilePayload struct {
	WalletAddress  string `json:"walletAddress"`
	HostIdentifier string `json:"hostIdentifier"`
}

// SenderProfilePayload is the payload for the sender profile endpoint
type SenderProfilePayload struct {
	WebhookURL      string          `json:"webhookURL"`
	DomainWhitelist []string        `json:"domainWhitelist"`
	FeePerTokenUnit decimal.Decimal `json:"feePerTokenUnit"`
	FeeAddress      string          `json:"feeAddress"`
	RefundAddress   string          `json:"refundAddress"`
}

// ProviderOrderTokenPayload defines the provider setting for a token
type ProviderOrderTokenPayload struct {
	Symbol                 string                                `json:"symbol" binding:"required"`
	ConversionRateType     providerordertoken.ConversionRateType `json:"conversionRateType" binding:"required"`
	FixedConversionRate    decimal.Decimal                       `json:"fixedConversionRate" binding:"required"`
	FloatingConversionRate decimal.Decimal                       `json:"floatingConversionRate" binding:"required"`
	MaxOrderAmount         decimal.Decimal                       `json:"maxOrderAmount" binding:"required"`
	MinOrderAmount         decimal.Decimal                       `json:"minOrderAmount" binding:"required"`
	Addresses              []struct {
		Address string `json:"address"`
		Network string `json:"network"`
	} `json:"addresses"`
}

// ProviderProfilePayload is the payload for the provider profile endpoint
type ProviderProfilePayload struct {
	TradingName          string                      `json:"tradingName"`
	Currency             string                      `json:"currency"`
	HostIdentifier       string                      `json:"hostIdentifier"`
	IsPartner            bool                        `json:"isPartner"`
	IsAvailable          bool                        `json:"isAvailable"`
	Tokens               []ProviderOrderTokenPayload `json:"tokens"`
	VisibilityMode       string                      `json:"visibilityMode"`
	Address              string                      `json:"address"`
	MobileNumber         string                      `json:"mobileNumber"`
	DateOfBirth          time.Time                   `json:"dateOfBirth"`
	BusinessName         string                      `json:"businessName"`
	IdentityDocumentType string                      `json:"identityType"`
	IdentityDocument     string                      `json:"identityDocument"`
	BusinessDocument     string                      `json:"businessDocument"`
}

// ProviderProfileResponse is the response for the provider profile endpoint
type ProviderProfileResponse struct {
	ID                   string                               `json:"id"`
	FirstName            string                               `json:"firstName"`
	LastName             string                               `json:"lastName"`
	Email                string                               `json:"email"`
	TradingName          string                               `json:"tradingName"`
	Currency             string                               `json:"currency"`
	Rate                 decimal.Decimal                      `json:"rate"`
	HostIdentifier       string                               `json:"hostIdentifier"`
	IsPartner            bool                                 `json:"isPartner"`
	IsAvailable          bool                                 `json:"isAvailable"`
	Tokens               []*ent.ProviderOrderToken            `json:"tokens"`
	APIKey               APIKeyResponse                       `json:"apiKey"`
	IsActive             bool                                 `json:"isActive"`
	Address              string                               `json:"address"`
	MobileNumber         string                               `json:"mobileNumber"`
	VisibilityMode       providerprofile.VisibilityMode       `json:"visibilityMode"`
	DateOfBirth          time.Time                            `json:"dateOfBirth"`
	BusinessName         string                               `json:"businessName"`
	IdentityDocumentType providerprofile.IdentityDocumentType `json:"identityType"`
	IdentityDocument     string                               `json:"identityDocument"`
	BusinessDocument     string                               `json:"businessDocument"`
	IsKybVerified        bool                                 `json:"isKybVerified"`
}

// ValidatorProfileResponse is the response for the validator profile endpoint
type ValidatorProfileResponse struct {
	ID             uuid.UUID      `json:"id"`
	WalletAddress  string         `json:"walletAddress"`
	HostIdentifier string         `json:"hostIdentifier"`
	APIKey         APIKeyResponse `json:"apiKey"`
}

// SenderProfileResponse is the response for the sender profile endpoint
type SenderProfileResponse struct {
	ID              uuid.UUID       `json:"id"`
	FirstName       string          `json:"firstName"`
	LastName        string          `json:"lastName"`
	Email           string          `json:"email"`
	WebhookURL      string          `json:"webhookUrl"`
	DomainWhitelist []string        `json:"domainWhitelist"`
	FeePerTokenUnit decimal.Decimal `json:"feePerTokenUnit"`
	FeeAddress      string          `json:"feeAddress"`
	RefundAddress   string          `json:"refundAddress"`
	APIKey          APIKeyResponse  `json:"apiKey"`
	IsActive        bool            `json:"isActive"`
}

// RefreshResponse is the response for the refresh endpoint
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

// APIKeyResponse is the response type for an API key
type APIKeyResponse struct {
	ID     uuid.UUID `json:"id"`
	Secret string    `json:"secret"`
}

// ERC20Transfer is the Transfer event of an ERC20 smart contract
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// LockPaymentOrderFields is the fields for a lock payment order
type LockPaymentOrderFields struct {
	ID                uuid.UUID
	Token             *ent.Token
	OrderID           string
	Amount            decimal.Decimal
	Rate              decimal.Decimal
	Label             string
	BlockNumber       int64
	TxHash            string
	Institution       string
	AccountIdentifier string
	AccountName       string
	ProviderID        string
	Memo              string
	ProvisionBucket   *ent.ProvisionBucket
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

// LockPaymentOrderResponse is the response for a lock payment order
type LockPaymentOrderResponse struct {
	ID                uuid.UUID       `json:"id"`
	Token             string          `json:"token"`
	OrderID           string          `json:"orderId"`
	Amount            decimal.Decimal `json:"amount"`
	Rate              decimal.Decimal `json:"rate"`
	Label             string          `json:"label"`
	BlockNumber       int64           `json:"blockNumber"`
	TxHash            string          `json:"txHash"`
	Institution       string          `json:"institution"`
	AccountIdentifier string          `json:"accountIdentifier"`
	AccountName       string          `json:"accountName"`
	ProviderID        string          `json:"providerId"`
	Memo              string          `json:"memo"`
	Network           string          `json:"network"`
	Status            string          `json:"status"`
	UpdatedAt         time.Time       `json:"updatedAt"`
	CreatedAt         time.Time       `json:"createdAt"`
}

// PaymentOrderRecipient describes a payment order recipient
type PaymentOrderRecipient struct {
	Institution       string `json:"institution" binding:"required"`
	AccountIdentifier string `json:"accountIdentifier" binding:"required"`
	AccountName       string `json:"accountName" binding:"required"`
	ProviderID        string `json:"providerId"`
	Memo              string `json:"memo" binding:"required"`
}

// NewPaymentOrderPayload is the payload for the create payment order endpoint
type NewPaymentOrderPayload struct {
	Amount          decimal.Decimal       `json:"amount" binding:"required"`
	Token           string                `json:"token" binding:"required"`
	Rate            decimal.Decimal       `json:"rate" binding:"required"`
	Network         string                `json:"network" binding:"required"`
	Recipient       PaymentOrderRecipient `json:"recipient" binding:"required"`
	Label           string                `json:"label" binding:"required"`
	FeePerTokenUnit decimal.Decimal       `json:"feePerTokenUnit"`
	FeeAddress      string                `json:"feeAddress"`
}

// ReceiveAddressResponse is the response type for a receive address
type ReceiveAddressResponse struct {
	ID             uuid.UUID       `json:"id"`
	Amount         decimal.Decimal `json:"amount"`
	Token          string          `json:"token"`
	Network        string          `json:"network"`
	ReceiveAddress string          `json:"receiveAddress"`
	ValidUntil     time.Time       `json:"validUntil"`
	SenderFee      decimal.Decimal `json:"senderFee"`
	TransactionFee decimal.Decimal `json:"transactionFee"`
}

// PaymentOrderResponse is the response type for a payment order
type PaymentOrderResponse struct {
	ID             uuid.UUID             `json:"id"`
	Amount         decimal.Decimal       `json:"amount"`
	AmountPaid     decimal.Decimal       `json:"amountPaid"`
	AmountReturned decimal.Decimal       `json:"amountReturned"`
	Token          string                `json:"token"`
	SenderFee      decimal.Decimal       `json:"senderFee"`
	TransactionFee decimal.Decimal       `json:"transactionFee"`
	Rate           decimal.Decimal       `json:"rate"`
	Network        string                `json:"network"`
	Label          string                `json:"label"`
	Recipient      PaymentOrderRecipient `json:"recipient"`
	FromAddress    string                `json:"fromAddress"`
	ReceiveAddress string                `json:"receiveAddress"`
	FeeAddress     string                `json:"feeAddress"`
	CreatedAt      time.Time             `json:"createdAt"`
	UpdatedAt      time.Time             `json:"updatedAt"`
	TxHash         string                `json:"txHash"`
	Status         paymentorder.Status   `json:"status"`
}

// PaymentOrderWebhookData is the data type for a payment order webhook
type PaymentOrderWebhookData struct {
	ID             uuid.UUID             `json:"id"`
	Amount         decimal.Decimal       `json:"amount"`
	AmountPaid     decimal.Decimal       `json:"amountPaid"`
	AmountReturned decimal.Decimal       `json:"amountReturned"`
	PercentSettled decimal.Decimal       `json:"percentSettled"`
	SenderFee      decimal.Decimal       `json:"senderFee"`
	NetworkFee     decimal.Decimal       `json:"networkFee"`
	Rate           decimal.Decimal       `json:"rate"`
	Network        string                `json:"network"`
	Label          string                `json:"label"`
	SenderID       uuid.UUID             `json:"senderId"`
	Recipient      PaymentOrderRecipient `json:"recipient"`
	FromAddress    string                `json:"fromAddress"`
	UpdatedAt      time.Time             `json:"updatedAt"`
	CreatedAt      time.Time             `json:"createdAt"`
	TxHash         string                `json:"txHash"`
	Status         paymentorder.Status   `json:"status"`
}

// PaymentOrderWebhookPayload is the request type for a payment order webhook
type PaymentOrderWebhookPayload struct {
	Event string                  `json:"event"`
	Data  PaymentOrderWebhookData `json:"data"`
}

// ConfirmEmailPayload is the payload for the confirmEmail endpoint
type ConfirmEmailPayload struct {
	Token string `json:"token" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// SendEmailPayload is content of a email request.
type SendEmailPayload struct {
	FromAddress string
	ToAddress   string
	Subject     string
	Body        string
	HTMLBody    string
}

// SendEmailResponse is the response for a sent email
type SendEmailResponse struct {
	Response string `json:"response"`
	Id       string `json:"id"`
}

// MarketRateResponse is the response for the market rate endpoint
type MarketRateResponse struct {
	MarketRate  decimal.Decimal `json:"marketRate"`
	MinimumRate decimal.Decimal `json:"minimumRate"`
	MaximumRate decimal.Decimal `json:"maximumRate"`
}

type ResendTokenPayload struct {
	Scope string `json:"scope" binding:"required,oneof=emailVerification resetPassword"`
	Email string `json:"email" binding:"required,email"`
}

type Institution struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Type string `json:"type"`
}

// SupportedCurrencies is the supported currencies response struct.
type SupportedCurrencies struct {
	Code       string          `json:"code"`
	Name       string          `json:"name"`
	ShortName  string          `json:"shortName"`
	Decimals   int8            `json:"decimals"`
	Symbol     string          `json:"symbol"`
	MarketRate decimal.Decimal `json:"marketRate"`
}

// Response is the struct for an API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorData is the struct for error data i.e when Status is "error"
type ErrorData struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Payload for reset password request
type ResetPasswordPayload struct {
	Password   string `json:"password" binding:"required,min=6,max=20"`
	ResetToken string `json:"resetToken"`
}

// Payload for reset password token endpoint
type ResetPasswordTokenPayload struct {
	Email string `json:"email"`
}

// ProviderLockOrderList is the struct for a list of provider lock orders
type ProviderLockOrderList struct {
	TotalRecords int                        `json:"total"`
	Page         int                        `json:"page"`
	PageSize     int                        `json:"pageSize"`
	Orders       []LockPaymentOrderResponse `json:"orders"`
}

// SenderOrderList is the struct for a list of sender payment orders
type SenderPaymentOrderList struct {
	TotalRecords int                    `json:"total"`
	Page         int                    `json:"page"`
	PageSize     int                    `json:"pageSize"`
	Orders       []PaymentOrderResponse `json:"orders"`
}

// ChangePasswordPayload is the payload for the change password endpoint
type ChangePasswordPayload struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6,max=20"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=20"`
}

// SenderStatsResponse is the response for the sender stats endpoint
type SenderStatsResponse struct {
	TotalOrders      int             `json:"totalOrders"`
	TotalOrderVolume decimal.Decimal `json:"totalOrderVolume"`
	TotalFeeEarnings decimal.Decimal `json:"totalFeeEarnings"`
}

// ProviderStatsResponse is the response for the provider stats endpoint
type ProviderStatsResponse struct {
	TotalOrders       int             `json:"totalOrders"`
	TotalFiatVolume   decimal.Decimal `json:"totalFiatVolume"`
	TotalCryptoVolume decimal.Decimal `json:"totalCryptoVolume"`
}
