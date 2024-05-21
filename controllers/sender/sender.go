package sender

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/storage"

	"github.com/paycrest/protocol/ent/network"
	"github.com/paycrest/protocol/ent/paymentorder"
	providerprofile "github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/senderprofile"
	"github.com/paycrest/protocol/ent/token"
	svc "github.com/paycrest/protocol/services"
	"github.com/paycrest/protocol/types"
	u "github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/logger"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// SenderController is a controller type for sender endpoints
type SenderController struct {
	indexerService        svc.Indexer
	receiveAddressService *svc.ReceiveAddressService
	orderService          svc.Order
}

// NewSenderController creates a new instance of SenderController
func NewSenderController(indexer svc.Indexer, order svc.Order) *SenderController {
	var indexerService svc.Indexer
	var orderService svc.Order

	if order == nil {
		orderService = svc.NewOrderService()
	} else {
		orderService = order
	}

	if indexer != nil {
		indexerService = indexer
	} else {
		indexerService = svc.NewIndexerService(orderService)
	}

	return &SenderController{
		indexerService:        indexerService,
		receiveAddressService: svc.NewReceiveAddressService(),
		orderService:          orderService,
	}
}

// InitiatePaymentOrder controller creates a payment order
func (ctrl *SenderController) InitiatePaymentOrder(ctx *gin.Context) {
	var payload types.NewPaymentOrderPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Get sender profile from the context
	senderCtx, ok := ctx.Get("sender")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	sender := senderCtx.(*ent.SenderProfile)

	conf := config.ServerConfig()

	if !sender.IsActive && !conf.Debug {
		u.APIResponse(ctx, http.StatusForbidden, "error", "Your account is not active", nil)
		return
	}

	// Get token from DB
	token, err := storage.Client.Token.
		Query().
		Where(
			token.SymbolEQ(payload.Token),
			token.HasNetworkWith(network.IdentifierEQ(payload.Network)),
			// TODO: check if token is enabled
		).
		WithNetwork().
		Only(ctx)
	if err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to validate payload", types.ErrorData{
			Field:   "Token",
			Message: "Provided token is not supported",
		})
		return
	}

	isPrivate := false
	isTokenNetworkPresent := false
	maxOrderAmount := decimal.NewFromInt(0)
	if payload.Recipients[0].ProviderID != "" {
		providerProfile, err := storage.Client.ProviderProfile.
			Query().
			Where(
				providerprofile.IDEQ(payload.Recipients[0].ProviderID),
			).
			WithOrderTokens().
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

	out:
		for _, orderToken := range providerProfile.Edges.OrderTokens {
			for _, address := range orderToken.Addresses {
				if address.Network == token.Edges.Network.Identifier {
					isTokenNetworkPresent = true
					break out
				}
			}
		}

		if !isTokenNetworkPresent {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "The selected network is not supported by the specified provider", nil)
			return
		}

		maxOrderAmount = providerProfile.Edges.OrderTokens[0].MaxOrderAmount

		if providerProfile.VisibilityMode == providerprofile.VisibilityModePrivate {
			isPrivate = true
		}
	}

	if isPrivate && payload.Amount.GreaterThan(maxOrderAmount) {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "The amount is beyond the maximum order amount for the specified provider", nil)
		return
	}
	// Generate receive address
	receiveAddress, err := ctrl.receiveAddressService.CreateSmartAccount(ctx, nil, nil)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to initiate payment order", nil)
		return
	}

	// Create payment order and recipient in a transaction
	tx, err := storage.Client.Tx(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to initiate payment order", nil)
		return
	}

	// Handle sender profile overrides
	var feePerTokenUnit decimal.Decimal
	var feeAddress string

	if payload.FeePerTokenUnit.IsZero() {
		feePerTokenUnit = sender.FeePerTokenUnit
	} else {
		feePerTokenUnit = payload.FeePerTokenUnit
	}

	if payload.FeeAddress == "" {
		feeAddress = sender.FeeAddress
	} else {
		if !sender.IsPartner {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to validate payload", types.ErrorData{
				Field:   "FeeAddress",
				Message: "FeeAddress is not allowed",
			})
			return
		}
		if !u.IsValidEthereumAddress(payload.FeeAddress) {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to validate payload", types.ErrorData{
				Field:   "FeeAddress",
				Message: "Invalid Ethereum address",
			})
			return
		}
		feeAddress = payload.FeeAddress
	}

	if payload.ReturnAddress != "" {
		if !u.IsValidEthereumAddress(payload.ReturnAddress) {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to validate payload", types.ErrorData{
				Field:   "ReturnAddress",
				Message: "Invalid Ethereum address",
			})
			return
		}
	}

	senderFee := feePerTokenUnit.Mul(payload.Amount).Div(payload.Rate).Round(int32(token.Decimals))
	protocolFee := payload.Amount.Mul(decimal.NewFromFloat(0.001)) // TODO: get protocol fee from contract -- currently 0.1%

	// Create payment order
	paymentOrder, err := tx.PaymentOrder.
		Create().
		SetSenderProfile(sender).
		SetAmount(payload.Amount).
		SetAmountPaid(decimal.NewFromInt(0)).
		SetAmountReturned(decimal.NewFromInt(0)).
		SetPercentSettled(decimal.NewFromInt(0)).
		SetNetworkFee(token.Edges.Network.Fee).
		SetProtocolFee(protocolFee).
		SetSenderFee(senderFee).
		SetToken(token).
		SetRate(payload.Rate).
		SetReceiveAddress(receiveAddress).
		SetReceiveAddressText(receiveAddress.Address).
		SetFeePerTokenUnit(feePerTokenUnit).
		SetFeeAddress(feeAddress).
		SetReturnAddress(payload.ReturnAddress).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to initiate payment order", nil)
		_ = tx.Rollback()
		return
	}

	// Create payment order recipient
	recipientsCreate := make([]*ent.PaymentOrderRecipientCreate, len(payload.Recipients))
	for i, recipient := range payload.Recipients {
		recipientsCreate[i] = tx.PaymentOrderRecipient.
			Create().
			SetInstitution(recipient.Institution).
			SetAccountIdentifier(recipient.AccountIdentifier).
			SetAccountName(recipient.AccountName).
			SetProviderID(recipient.ProviderID).
			SetMemo(recipient.Memo).
			SetPaymentOrder(paymentOrder)
	}

	_, err = tx.PaymentOrderRecipient.CreateBulk(recipientsCreate...).Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to initiate payment order", nil)
		_ = tx.Rollback()
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to initiate payment order", nil)
		return
	}

	// Start a background process to index token transfers through the lifecycle of the receive address
	go func() {
		_ = ctrl.indexerService.IndexERC20Transfer(ctx, nil, receiveAddress)
	}()

	u.APIResponse(ctx, http.StatusCreated, "success", "Payment order initiated successfully",
		&types.ReceiveAddressResponse{
			ID:             paymentOrder.ID,
			Amount:         paymentOrder.Amount,
			Token:          payload.Token,
			Network:        token.Edges.Network.Identifier,
			ReceiveAddress: receiveAddress.Address,
			ValidUntil:     receiveAddress.ValidUntil,
			SenderFee:      senderFee,
			TransactionFee: protocolFee.Add(token.Edges.Network.Fee),
		})
}

// GetPaymentOrderByID controller fetches a payment order by ID
func (ctrl *SenderController) GetPaymentOrderByID(ctx *gin.Context) {
	// Get order ID from the URL
	orderID := ctx.Param("id")

	// Convert order ID to UUID
	id, err := uuid.Parse(orderID)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Invalid order ID", nil)
		return
	}

	// Get sender profile from the context
	senderCtx, ok := ctx.Get("sender")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	sender := senderCtx.(*ent.SenderProfile)

	// Fetch payment order from the database
	paymentOrder, err := storage.Client.PaymentOrder.
		Query().
		Where(
			paymentorder.IDEQ(id),
			paymentorder.HasSenderProfileWith(senderprofile.IDEQ(sender.ID)),
		).
		WithRecipient().
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		Only(ctx)

	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error",
			"Payment order not found", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "The order has been successfully retrieved", &types.PaymentOrderResponse{
		ID:             paymentOrder.ID,
		Amount:         paymentOrder.Amount,
		AmountPaid:     paymentOrder.AmountPaid,
		AmountReturned: paymentOrder.AmountReturned,
		Token:          paymentOrder.Edges.Token.Symbol,
		SenderFee:      paymentOrder.SenderFee,
		TransactionFee: paymentOrder.NetworkFee.Add(paymentOrder.ProtocolFee),
		Rate:           paymentOrder.Rate,
		Network:        paymentOrder.Edges.Token.Edges.Network.Identifier,
		Recipient: types.PaymentOrderRecipient{
			Institution:       paymentOrder.Edges.Recipient.Institution,
			AccountIdentifier: paymentOrder.Edges.Recipient.AccountIdentifier,
			AccountName:       paymentOrder.Edges.Recipient.AccountName,
			ProviderID:        paymentOrder.Edges.Recipient.ProviderID,
			Memo:              paymentOrder.Edges.Recipient.Memo,
		},
		FromAddress:    paymentOrder.FromAddress,
		ReturnAddress:  paymentOrder.ReturnAddress,
		ReceiveAddress: paymentOrder.ReceiveAddressText,
		FeeAddress:     paymentOrder.FeeAddress,
		GatewayID:      paymentOrder.GatewayID,
		CreatedAt:      paymentOrder.CreatedAt,
		UpdatedAt:      paymentOrder.UpdatedAt,
		TxHash:         paymentOrder.TxHash,
		Status:         paymentOrder.Status,
	})
}

// GetPaymentOrders controller fetches all payment orders
func (ctrl *SenderController) GetPaymentOrders(ctx *gin.Context) {
	// Get sender profile from the context
	senderCtx, ok := ctx.Get("sender")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	sender := senderCtx.(*ent.SenderProfile)

	// Get ordering query param
	ordering := ctx.Query("ordering")
	order := ent.Desc(paymentorder.FieldCreatedAt)
	if ordering == "asc" {
		order = ent.Asc(paymentorder.FieldCreatedAt)
	}

	// Get page and pageSize query params
	page, offset, pageSize := u.Paginate(ctx)

	paymentOrderQuery := storage.Client.PaymentOrder.Query()

	// Filter by sender
	paymentOrderQuery = paymentOrderQuery.Where(
		paymentorder.HasSenderProfileWith(senderprofile.IDEQ(sender.ID)),
	)

	// Filter by status
	statusQueryParam := ctx.Query("status")
	statusMap := map[string]paymentorder.Status{
		"initiated": paymentorder.StatusInitiated,
		"pending":   paymentorder.StatusPending,
		"reverted":  paymentorder.StatusReverted,
		"expired":   paymentorder.StatusExpired,
		"settled":   paymentorder.StatusSettled,
		"refunded":  paymentorder.StatusRefunded,
	}

	if status, ok := statusMap[statusQueryParam]; ok {
		paymentOrderQuery = paymentOrderQuery.Where(
			paymentorder.StatusEQ(status),
		)
	}

	// Filter by token
	tokenQueryParam := ctx.Query("token")

	if tokenQueryParam != "" {
		tokenExists, err := storage.Client.Token.
			Query().
			Where(
				token.SymbolEQ(tokenQueryParam),
			).
			Exist(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to fetch payment orders", nil)
			return
		}

		if tokenExists {
			paymentOrderQuery = paymentOrderQuery.Where(
				paymentorder.HasTokenWith(
					token.SymbolEQ(tokenQueryParam),
				),
			)
		}
	}

	// Filter by network
	networkQueryParam := ctx.Query("network")

	if networkQueryParam != "" {
		networkExists, err := storage.Client.Network.
			Query().
			Where(
				network.IdentifierEQ(networkQueryParam),
			).
			Exist(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to fetch payment orders", nil)
			return
		}

		if networkExists {
			paymentOrderQuery = paymentOrderQuery.Where(
				paymentorder.HasTokenWith(
					token.HasNetworkWith(
						network.IdentifierEQ(networkQueryParam),
					),
				),
			)
		}
	}

	count, err := paymentOrderQuery.Count(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch payment orders", nil)
		return
	}

	// Fetch payment orders
	paymentOrders, err := paymentOrderQuery.
		WithRecipient().
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		Limit(pageSize).
		Offset(offset).
		Order(order).
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to fetch payment orders", nil)
		return
	}

	var orders []types.PaymentOrderResponse

	for _, paymentOrder := range paymentOrders {
		orders = append(orders, types.PaymentOrderResponse{
			ID:             paymentOrder.ID,
			Amount:         paymentOrder.Amount,
			AmountPaid:     paymentOrder.AmountPaid,
			AmountReturned: paymentOrder.AmountReturned,
			Token:          paymentOrder.Edges.Token.Symbol,
			SenderFee:      paymentOrder.SenderFee,
			TransactionFee: paymentOrder.NetworkFee.Add(paymentOrder.ProtocolFee),
			Rate:           paymentOrder.Rate,
			Network:        paymentOrder.Edges.Token.Edges.Network.Identifier,
			Recipient: types.PaymentOrderRecipient{
				Institution:       paymentOrder.Edges.Recipient.Institution,
				AccountIdentifier: paymentOrder.Edges.Recipient.AccountIdentifier,
				AccountName:       paymentOrder.Edges.Recipient.AccountName,
				ProviderID:        paymentOrder.Edges.Recipient.ProviderID,
				Memo:              paymentOrder.Edges.Recipient.Memo,
			},
			FromAddress:    paymentOrder.FromAddress,
			ReturnAddress:  paymentOrder.ReturnAddress,
			ReceiveAddress: paymentOrder.ReceiveAddressText,
			FeeAddress:     paymentOrder.FeeAddress,
			GatewayID:      paymentOrder.GatewayID,
			CreatedAt:      paymentOrder.CreatedAt,
			UpdatedAt:      paymentOrder.UpdatedAt,
			TxHash:         paymentOrder.TxHash,
			Status:         paymentOrder.Status,
		})
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Payment orders retrieved successfully", types.SenderPaymentOrderList{
		Page:         page,
		PageSize:     pageSize,
		TotalRecords: count,
		Orders:       orders,
	})
}

// Stats controller fetches sender stats
func (ctrl *SenderController) Stats(ctx *gin.Context) {
	// Get sender profile from the context
	senderCtx, ok := ctx.Get("sender")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	sender := senderCtx.(*ent.SenderProfile)

	// Aggregate sender stats from db

	var w []struct {
		Sum decimal.Decimal
	}
	err := storage.Client.PaymentOrder.
		Query().
		Where(paymentorder.HasSenderProfileWith(senderprofile.IDEQ(sender.ID))).
		Aggregate(
			ent.Sum(paymentorder.FieldAmount),
		).
		Scan(ctx, &w)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch sender stats", nil)
		return
	}

	var v []struct {
		Count int
		Sum   decimal.Decimal
	}
	err = storage.Client.PaymentOrder.
		Query().
		Where(paymentorder.HasSenderProfileWith(senderprofile.IDEQ(sender.ID))).
		Aggregate(
			ent.Count(),
			ent.Sum(paymentorder.FieldSenderFee),
		).
		Scan(ctx, &v)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch sender stats", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Sender stats retrieved successfully", types.SenderStatsResponse{
		TotalOrders:      v[0].Count,
		TotalOrderVolume: w[0].Sum,
		TotalFeeEarnings: v[0].Sum,
	})
}
