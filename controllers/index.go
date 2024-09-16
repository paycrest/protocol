package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/ent/identityverificationrequest"
	"github.com/paycrest/protocol/ent/institution"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/token"
	svc "github.com/paycrest/protocol/services"
	orderSvc "github.com/paycrest/protocol/services/order"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	u "github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/logger"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

var cryptoConf = config.CryptoConfig()
var serverConf = config.ServerConfig()
var identityConf = config.IdentityConfig()

// Controller is the default controller for other endpoints
type Controller struct {
	orderService         types.OrderService
	priorityQueueService *svc.PriorityQueueService
}

// NewController creates a new instance of AuthController with injected services
func NewController() *Controller {
	return &Controller{
		orderService:         orderSvc.NewOrderEVM(),
		priorityQueueService: svc.NewPriorityQueueService(),
	}
}

// GetFiatCurrencies controller fetches the supported fiat currencies
func (ctrl *Controller) GetFiatCurrencies(ctx *gin.Context) {
	// fetch stored fiat currencies.
	fiatcurrencies, err := storage.Client.FiatCurrency.
		Query().
		Where(fiatcurrency.IsEnabledEQ(true)).
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to fetch FiatCurrencies", err.Error())
		return
	}

	currencies := make([]types.SupportedCurrencies, 0, len(fiatcurrencies))
	for _, currency := range fiatcurrencies {
		currencies = append(currencies, types.SupportedCurrencies{
			Code:       currency.Code,
			Name:       currency.Name,
			ShortName:  currency.ShortName,
			Decimals:   int8(currency.Decimals),
			Symbol:     currency.Symbol,
			MarketRate: currency.MarketRate,
		})
	}

	u.APIResponse(ctx, http.StatusOK, "success", "OK", currencies)
}

// GetInstitutionsByCurrency controller fetches the supported institutions for a given currency
func (ctrl *Controller) GetInstitutionsByCurrency(ctx *gin.Context) {
	// Get currency code from the URL
	currencyCode := ctx.Param("currency_code")

	institutions, err := storage.Client.Institution.
		Query().
		Where(institution.HasFiatCurrencyWith(
			fiatcurrency.CodeEQ(currencyCode),
		)).
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to fetch institutions", nil)
		return
	}

	response := make([]types.SupportedInstitutions, 0, len(institutions))
	for _, institution := range institutions {
		response = append(response, types.SupportedInstitutions{
			Code: institution.Code,
			Name: institution.Name,
			Type: institution.Type,
		})
	}

	u.APIResponse(ctx, http.StatusOK, "success", "OK", response)
}

// GetTokenRate controller fetches the current rate of the cryptocurrency token against the fiat currency
func (ctrl *Controller) GetTokenRate(ctx *gin.Context) {
	// Parse path parameters
	token, err := storage.Client.Token.
		Query().
		Where(
			token.SymbolEQ(strings.ToUpper(ctx.Param("token"))),
			token.IsEnabledEQ(true),
		).
		First(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch token rate", nil)
		return
	}

	if token == nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Token is not supported", nil)
		return
	}

	currency, err := storage.Client.FiatCurrency.
		Query().
		Where(
			fiatcurrency.IsEnabledEQ(true),
			fiatcurrency.CodeEQ(strings.ToUpper(ctx.Param("fiat"))),
		).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Fiat currency is not supported", nil)
		return
	}

	tokenAmount, err := decimal.NewFromString(ctx.Param("amount"))
	if err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid amount", nil)
		return
	}

	rateResponse := decimal.NewFromInt(0)

	// get providerID from query params
	providerID := ctx.Query("provider_id")
	if providerID != "" {
		// get the provider from the bucket
		provider, err := storage.Client.ProviderProfile.
			Query().
			Where(providerprofile.IDEQ(providerID)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				u.APIResponse(ctx, http.StatusBadRequest, "error", "Provider not found", nil)
				return
			} else {
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch provider profile", nil)
				return
			}
		}

		rateResponse, err = ctrl.priorityQueueService.GetProviderRate(ctx, provider, token.Symbol)
		if err != nil {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch provider rate", nil)
			return
		}

	} else {
		// Get redis keys for provision buckets
		keys, _, err := storage.RedisClient.Scan(ctx, uint64(0), "bucket_"+currency.Code+"_*_*", 100).Result()
		if err != nil {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
			return
		}

		highestMaxAmount := decimal.NewFromInt(0)

		// Scan through the buckets to find a matching rate
		for _, key := range keys {
			bucketData := strings.Split(key, "_")
			minAmount, _ := decimal.NewFromString(bucketData[2])
			maxAmount, _ := decimal.NewFromString(bucketData[3])

			// Get the topmost provider in the priority queue of the bucket
			providerData, err := storage.RedisClient.LIndex(ctx, key, 0).Result()
			if err != nil {
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
				return
			}

			// Get fiat equivalent of the token amount
			rate, _ := decimal.NewFromString(strings.Split(providerData, ":")[2])
			fiatAmount := tokenAmount.Mul(rate)

			// Check if fiat amount is within the bucket range and set the rate
			if fiatAmount.GreaterThanOrEqual(minAmount) && fiatAmount.LessThanOrEqual(maxAmount) {
				rateResponse = rate
				break
			} else {
				// Get the highest max amount
				if maxAmount.GreaterThan(highestMaxAmount) {
					highestMaxAmount = maxAmount
					rateResponse = rate
				}
			}
		}
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Rate fetched successfully", rateResponse)
}

// GetAggregatorPublicKey controller expose Aggregator Public Key
func (ctrl *Controller) GetAggregatorPublicKey(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", cryptoConf.AggregatorPublicKey)
}

// VerifyAccount controller verifies an account of a given institution
func (ctrl *Controller) VerifyAccount(ctx *gin.Context) {
	var payload types.VerifyAccountRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	institution, err := storage.Client.Institution.
		Query().
		Where(institution.CodeEQ(payload.Institution)).
		WithFiatCurrency().
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to validate payload", []types.ErrorData{{
			Field:   "Institution",
			Message: "Institution is not supported",
		}})
		return
	}

	// TODO: Remove this after testing non-NGN institutions
	if institution.Edges.FiatCurrency.Code != "NGN" {
		u.APIResponse(ctx, http.StatusOK, "success", "Account name was fetched successfully", "OK")
		return
	}

	providers, err := storage.Client.ProviderProfile.
		Query().
		Where(
			providerprofile.HasCurrencyWith(
				fiatcurrency.CodeEQ(institution.Edges.FiatCurrency.Code),
			),
			providerprofile.HostIdentifierNotNil(),
			providerprofile.IsActiveEQ(true),
			providerprofile.IsAvailableEQ(true),
		).
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to verify account", err.Error())
		return
	}

	var res fastshot.Response
	var data map[string]interface{}
	for _, provider := range providers {
		res, err = fastshot.NewClient(provider.HostIdentifier).
			Config().SetTimeout(30 * time.Second).
			Build().POST("/verify_account").
			Body().AsJSON(payload).
			Send()
		if err != nil {
			continue
		}

		data, err = u.ParseJSONResponse(res.RawResponse)
		if err != nil {
			continue
		}
	}

	if err != nil {
		logger.Errorf("error: %v %v", err, data)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to verify account", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Account name was fetched successfully", data["data"].(string))
}

// GetLockPaymentOrderStatus controller fetches a payment order status by ID
func (ctrl *Controller) GetLockPaymentOrderStatus(ctx *gin.Context) {
	// Get order ID from the URL
	orderID := ctx.Param("id")

	// Fetch related payment orders from the database
	orders, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.GatewayIDEQ(orderID),
		).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		WithTransactions().
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch order status", nil)
		return
	}

	var settlements []types.LockPaymentOrderSplitOrder
	var receipts []types.LockPaymentOrderTxReceipt
	var settlePercent decimal.Decimal
	var totalAmount decimal.Decimal

	for _, order := range orders {
		for _, transaction := range order.Edges.Transactions {
			if u.ContainsString([]string{"order_settled", "order_created", "order_refunded"}, transaction.Status.String()) {
				var status lockpaymentorder.Status
				if transaction.Status.String() == "order_created" {
					status = lockpaymentorder.StatusPending
				} else {
					status = lockpaymentorder.Status(strings.TrimPrefix(transaction.Status.String(), "order_"))
				}
				receipts = append(receipts, types.LockPaymentOrderTxReceipt{
					Status:    status,
					TxHash:    transaction.TxHash,
					Timestamp: transaction.CreatedAt,
				})
			}
		}

		settlements = append(settlements, types.LockPaymentOrderSplitOrder{
			SplitOrderID: order.ID,
			Amount:       order.Amount,
			Rate:         order.Rate,
			OrderPercent: order.OrderPercent,
		})

		settlePercent = settlePercent.Add(order.OrderPercent)
		totalAmount = totalAmount.Add(order.Amount)
	}

	// Sort receipts by latest timestamp
	slices.SortStableFunc(receipts, func(a, b types.LockPaymentOrderTxReceipt) int {
		return b.Timestamp.Compare(a.Timestamp)
	})

	if (len(orders) == 0) || (len(receipts) == 0) {
		u.APIResponse(ctx, http.StatusNotFound, "error", "Order not found", nil)
		return
	}

	response := &types.LockPaymentOrderStatusResponse{
		OrderID:       orders[0].GatewayID,
		Amount:        totalAmount,
		Token:         orders[0].Edges.Token.Symbol,
		Network:       orders[0].Edges.Token.Edges.Network.Identifier,
		SettlePercent: settlePercent,
		Status:        orders[0].Status,
		TxHash:        receipts[0].TxHash,
		Settlements:   settlements,
		TxReceipts:    receipts,
		UpdatedAt:     orders[0].UpdatedAt,
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Order status fetched successfully", response)
}

// RequestIDVerification controller requests identity verification details
func (ctrl *Controller) RequestIDVerification(ctx *gin.Context) {
	var payload types.NewIDVerificationRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Validate wallet signature
	signature, err := hex.DecodeString(payload.Signature)
	if err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid signature", "Signature is not in the correct format")
		return
	}

	if len(signature) != 65 {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid signature", "Signature length is not correct")
		return
	}

	if signature[64] != 27 && signature[64] != 28 {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid signature", "Invalid recovery ID")
		return
	}
	signature[64] -= 27

	// Verify wallet signature
	message := fmt.Sprintf("I accept the KYC Policy and hereby request an identity verification check for %s with nonce %s", payload.WalletAddress, payload.Nonce)

	prefix := "\x19Ethereum Signed Message:\n" + fmt.Sprint(len(message))
	hash := crypto.Keccak256Hash([]byte(prefix + message))

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid signature", nil)
		return
	}

	recoveredAddress := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	if !strings.EqualFold(recoveredAddress.Hex(), payload.WalletAddress) {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid signature", nil)
		return
	}

	// Check if there is an existing verification request
	ivr, err := storage.Client.IdentityVerificationRequest.
		Query().
		Where(
			identityverificationrequest.WalletAddressEQ(payload.WalletAddress),
		).
		Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to request identity verification", nil)
			return
		}
	}

	timestamp := time.Now()
	if ivr.Status == "pending" && ivr.LastURLCreatedAt.Add(24*7*time.Hour).Before(timestamp) {
		// Request is expired, delete db entry
		_, err := storage.Client.IdentityVerificationRequest.
			Delete().
			Where(
				identityverificationrequest.WalletAddressEQ(payload.WalletAddress),
			).
			Exec(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to request identity verification", nil)
			return
		}

	} else if ivr.Status == "pending" && (ivr.LastURLCreatedAt.Add(24*7*time.Hour).Equal(timestamp) || ivr.LastURLCreatedAt.Add(24*7*time.Hour).After(timestamp)) {
		// Update the wallet signature in db
		_, err = storage.Client.IdentityVerificationRequest.
			Update().
			SetWalletSignature(payload.Signature).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to request identity verification", nil)
			return
		}

		u.APIResponse(ctx, http.StatusOK, "success", "This account has a pending identity verification", &types.NewIDVerificationResponse{
			URL:       ivr.VerificationURL,
			ExpiresAt: ivr.LastURLCreatedAt,
		})
		return
	}

	if ivr.Status == "success" {
		u.APIResponse(ctx, http.StatusBadRequest, "success", "Failed to request identity verification", "This account has already been successfully verified")
		return
	}

	// Generate Smile Identity signature
	h := hmac.New(sha256.New, []byte(identityConf.SmileIdentityApiKey))
	h.Write([]byte(timestamp.Format(time.RFC3339)))
	h.Write([]byte(identityConf.SmileIdentityPartnerId))
	h.Write([]byte("sid_request"))

	// Initiate KYC verification
	res, err := fastshot.NewClient(identityConf.SmileIdentityBaseUrl).
		Config().SetTimeout(30 * time.Second).
		Build().POST("/v1/smile_links").
		Body().AsJSON(map[string]interface{}{
		"partner_id":   identityConf.SmileIdentityPartnerId,
		"signature":    base64.StdEncoding.EncodeToString(h.Sum(nil)),
		"timestamp":    timestamp,
		"name":         "Aggregator KYC",
		"company_name": "Paycrest Aggregator",
		"id_types": []map[string]interface{}{
			// Nigeria
			{
				"country":             "NG",
				"id_type":             "PASSPORT",
				"verification_method": "doc_verification",
			},
			{
				"country":             "NG",
				"id_type":             "DRIVERS_LICENSE",
				"verification_method": "doc_verification",
			},
			{
				"country":             "NG",
				"id_type":             "V_NIN",
				"verification_method": "biometric_kyc",
			},
			{
				"country":             "NG",
				"id_type":             "VOTER_ID",
				"verification_method": "biometric_kyc",
			},
			{
				"country":             "NG",
				"id_type":             "RESIDENT_ID",
				"verification_method": "doc_verification",
			},
			{
				"country":             "NG",
				"id_type":             "IDENTITY_CARD",
				"verification_method": "doc_verification",
			},

			// Ghana
			{
				"country":             "GH",
				"id_type":             "PASSPORT",
				"verification_method": "enhanced_document_verification",
			},
			{
				"country":             "GH",
				"id_type":             "VOTER_ID",
				"verification_method": "enhanced_document_verification",
			},
			{
				"country":             "GH",
				"id_type":             "NEW_VOTER_ID",
				"verification_method": "biometric_kyc",
			},
			{
				"country":             "GH",
				"id_type":             "DRIVERS_LICENSE",
				"verification_method": "doc_verification",
			},
			{
				"country":             "GH",
				"id_type":             "SSNIT",
				"verification_method": "biometric_kyc",
			},

			// Kenya
			{
				"country":             "KE",
				"id_type":             "PASSPORT",
				"verification_method": "enhanced_document_verification",
			},
			{
				"country":             "KE",
				"id_type":             "DRIVERS_LICENSE",
				"verification_method": "doc_verification",
			},
			{
				"country":             "KE",
				"id_type":             "ALIEN_CARD",
				"verification_method": "biometric_kyc",
			},
			{
				"country":             "KE",
				"id_type":             "NATIONAL_ID",
				"verification_method": "biometric_kyc",
			},

			// South Africa
			// {
			// 	"country":             "ZA",
			// 	"id_type":             "PASSPORT",
			// 	"verification_method": "doc_verification",
			// },
			// {
			// 	"country":             "ZA",
			// 	"id_type":             "DRIVERS_LICENSE",
			// 	"verification_method": "doc_verification",
			// },
			// {
			// 	"country":             "ZA",
			// 	"id_type":             "RESIDENT_ID",
			// 	"verification_method": "doc_verification",
			// },
			// {
			// 	"country":             "ZA",
			// 	"id_type":             "NATIONAL_ID",
			// 	"verification_method": "biometric_kyc",
			// },
		},
		"callback_url":            fmt.Sprintf("%s/v1/kyc/webhook", serverConf.HostDomain),
		"data_privacy_policy_url": "https://www.paycrest.io/privacy-policy",
		"logo_url":                "https://i.postimg.cc/Twrq0gjC/mark-2x-2.png",
		"is_single_use":           true,
		"user_id":                 payload.WalletAddress,
		"expires_at":              timestamp.Add(24 * time.Hour).Format(time.RFC3339Nano),
	}).
		Send()
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to request identity verification", "Couldn't reach identity provider")
		return
	}

	data, err := u.ParseJSONResponse(res.RawResponse)
	if err != nil {
		logger.Errorf("error: %v %v", err, data)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to request identity verification", data)
		return
	}

	// Save the verification details to the database
	ivr, err = storage.Client.IdentityVerificationRequest.
		Create().
		SetWalletAddress(payload.WalletAddress).
		SetPlatform("smile_id").
		SetPlatformRef(data["ref_id"].(string)).
		SetVerificationURL(data["link"].(string)).
		SetLastURLCreatedAt(timestamp.Add(24 * time.Hour)).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to request identity verification", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Identity verification requested successfully", &types.NewIDVerificationResponse{
		URL:       ivr.VerificationURL,
		ExpiresAt: ivr.LastURLCreatedAt,
	})
}

// GetIDVerificationStatus controller fetches the status of an identity verification request
func (ctrl *Controller) GetIDVerificationStatus(ctx *gin.Context) {
	// Get wallet address from the URL
	walletAddress := ctx.Param("wallet_address")

	// Fetch related identity verification request from the database
	ivr, err := storage.Client.IdentityVerificationRequest.
		Query().
		Where(
			identityverificationrequest.WalletAddressEQ(walletAddress),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// Check the platform's status endpoint
			u.APIResponse(ctx, http.StatusNotFound, "error", "No verification request found for this wallet address", nil)
			return
		}
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch identity verification status", nil)
		return
	}

	response := struct {
		Status string `json:"status"`
	}{
		Status: ivr.Status.String(),
	}

	// Check if the status is pending and fetch the status from the platform
	if ivr.Status == identityverificationrequest.StatusPending {
		status, err := getSmileLinkStatus(ivr.PlatformRef)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to fetch identity verification status", nil)
			return
		}

		response.Status = status

		if status != "pending" {
			// Update the verification status in the database if it's not pending
			_, err := storage.Client.IdentityVerificationRequest.
				Update().
				Where(
					identityverificationrequest.WalletAddressEQ(walletAddress),
				).
				SetStatus(identityverificationrequest.Status(response.Status)).
				Save(ctx)
			if err != nil {
				logger.Errorf("error: %v", err)
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update identity verification status", nil)
				return
			}
		}
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Identity verification status fetched successfully", response)
}

// KYCWebhook handles the webhook callback from Smile Identity
func (ctrl *Controller) KYCWebhook(ctx *gin.Context) {
	var payload types.SmileIDWebhookPayload

	// Parse the JSON payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("Failed to parse webhook payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Verify the webhook signature
	if !verifySmileIDWebhookSignature(payload, payload.Signature) {
		logger.Errorf("Invalid webhook signature")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// Process the webhook
	status := identityverificationrequest.StatusPending

	// Check for success codes
	successCodes := []string{
		"0810", // Document Verified
		"1020", // Exact Match (Basic KYC and Enhanced KYC)
		"1012", // Valid ID / ID Number Validated (Enhanced KYC)
		"0820", // Authenticate User Machine Judgement - PASS
		"0840", // Enroll User PASS - Machine Judgement
	}

	// Check for failed codes
	failedCodes := []string{
		"0811", // Document Not Verified
		"0812", // Document Not Verified
		"0813", // Document Not Verified - Machine Judgement
		"1022", // No Match
		"1023", // No Found
		"1011", // Invalid ID / ID Number Invalid
		"1013", // ID Number Not Found
		"1014", // Unsupported ID Type
		"0821", // Images did not match
		"0911", // No Face Found
		"0912", // Face Not Matching
		"0921", // Face Not Found
		"0922", // Selfie Quality Too Poor
		"0841", // Enroll User FAIL
		"0941", // Face Not Found
		"0942", // Face Poor Quality
	}

	if slices.Contains(successCodes, payload.ResultCode) {
		status = identityverificationrequest.StatusSuccess
	}

	if slices.Contains(failedCodes, payload.ResultCode) {
		status = identityverificationrequest.StatusFailed
	}

	// Update the verification status in the database
	_, err := storage.Client.IdentityVerificationRequest.
		Update().
		Where(
			identityverificationrequest.WalletAddressEQ(payload.PartnerParams.UserID),
			identityverificationrequest.StatusEQ(identityverificationrequest.StatusPending),
		).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		logger.Errorf("Failed to update verification status: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

// verifyWebhookSignature verifies the signature of a Smile Identity webhook
func verifySmileIDWebhookSignature(payload types.SmileIDWebhookPayload, receivedSignature string) bool {
	// Create HMAC
	// Generate Smile Identity signature
	h := hmac.New(sha256.New, []byte(identityConf.SmileIdentityApiKey))
	h.Write([]byte(payload.Timestamp))
	h.Write([]byte(identityConf.SmileIdentityPartnerId))
	h.Write([]byte("sid_request"))

	// Compare the computed signature with the one in the header
	computedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return computedSignature == receivedSignature
}

// getSmileLinkStatus fetches the status of a Smile Link
func getSmileLinkStatus(linkID string) (string, error) {
	// Generate signature
	timestamp := time.Now().UTC().Format(time.RFC3339)
	h := hmac.New(sha256.New, []byte(identityConf.SmileIdentityApiKey))
	h.Write([]byte(timestamp))
	h.Write([]byte(identityConf.SmileIdentityPartnerId))
	h.Write([]byte("sid_request"))

	// Get Smile Link status
	res, err := fastshot.NewClient(identityConf.SmileIdentityBaseUrl).
		Config().SetTimeout(30 * time.Second).
		Build().POST(fmt.Sprintf("/v1/smile_links/%s", linkID)).
		Body().AsJSON(map[string]interface{}{
		"partner_id": identityConf.SmileIdentityPartnerId,
		"signature":  base64.StdEncoding.EncodeToString(h.Sum(nil)),
		"timestamp":  timestamp,
	}).
		Send()
	if err != nil {
		return "", fmt.Errorf("failed to get Smile Link status: %w", err)
	}

	data, err := u.ParseJSONResponse(res.RawResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse Smile Link response: %w", err)
	}

	totalJobs := data["total_jobs"].(float64)
	successfulJobs := data["successful_jobs"].(float64)
	failedJobs := data["failed_jobs"].(float64)

	if failedJobs+successfulJobs < totalJobs {
		return "pending", nil
	}

	if totalJobs == successfulJobs && failedJobs == 0 {
		return "success", nil
	}

	if failedJobs > 0 {
		return "failed", nil
	}

	return "", nil
}
