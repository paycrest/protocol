package sender

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent"
	db "github.com/paycrest/paycrest-protocol/storage"

	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/ent/token"
	svc "github.com/paycrest/paycrest-protocol/services"
	"github.com/paycrest/paycrest-protocol/types"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// SenderController is a controller type for sender endpoints
type SenderController struct {
	indexerService        *svc.IndexerService
	receiveAddressService *svc.ReceiveAddressService
	orderService          *svc.OrderService
}

// NewSenderController creates a new instance of SenderController
func NewSenderController(indexer svc.Indexer) *SenderController {
	var indexerService *svc.IndexerService

	if indexer != nil {
		indexerService = svc.NewIndexerService(indexer)
	} else {
		indexerService = svc.NewIndexerService(nil)
	}

	return &SenderController{
		indexerService:        indexerService,
		receiveAddressService: svc.NewReceiveAddressService(),
		orderService:          svc.NewOrderService(),
	}
}

// CreatePaymentOrder controller creates a payment order
func (ctrl *SenderController) CreatePaymentOrder(ctx *gin.Context) {
	var payload types.NewPaymentOrderPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Generate receive address
	receiveAddress, err := ctrl.receiveAddressService.CreateSmartAccount(ctx, nil, nil)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", err.Error())
		return
	}

	// Get token from DB
	token, err := db.Client.Token.
		Query().
		Where(token.SymbolEQ(payload.Token)).
		WithNetwork().
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Provided crypto token is not supported", err.Error())
		return
	}

	// Create payment order and recipient in a transaction
	tx, err := db.Client.Tx(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", err.Error())
		return
	}

	// Create payment order
	apiKey, _ := ctx.Get("api_key")

	paymentOrder, err := tx.PaymentOrder.
		Create().
		SetAPIKey(apiKey.(*ent.APIKey)).
		SetAmount(payload.Amount).
		SetAmountPaid(decimal.NewFromInt(0)).
		SetToken(token).
		SetReceiveAddress(receiveAddress).
		SetReceiveAddressText(receiveAddress.Address).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", err.Error())
		_ = tx.Rollback()
		return
	}

	// Create payment order recipient
	_, err = tx.PaymentOrderRecipient.
		Create().
		SetInstitution(payload.Recipient.Institution).
		SetAccountIdentifier(payload.Recipient.AccountIdentifier).
		SetAccountName(payload.Recipient.AccountName).
		SetProviderID(payload.Recipient.ProviderID).
		SetPaymentOrder(paymentOrder).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", err.Error())
		_ = tx.Rollback()
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", err.Error())
		return
	}

	// Start a background process to index token transfers to the receive address
	go ctrl.indexerService.RunIndexERC20Transfer(ctx, receiveAddress)

	paymentOrderAmount, _ := paymentOrder.Amount.Float64()

	u.APIResponse(ctx, http.StatusCreated, "success", "Payment order initiated successfully",
		&types.ReceiveAddressResponse{
			ID:             paymentOrder.ID,
			Amount:         paymentOrderAmount,
			Network:        token.Edges.Network.Identifier.String(),
			ReceiveAddress: paymentOrder.ReceiveAddressText,
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
			"Invalid order ID", err.Error())
		return
	}

	// Fetch payment order from the database
	paymentOrder, err := db.Client.PaymentOrder.
		Query().
		Where(paymentorder.ID(id)).
		WithRecipient().
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		First(ctx)

	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error",
			"Payment order not found", err.Error())
		return
	}

	paymentOrderAmount, _ := paymentOrder.Amount.Float64()

	u.APIResponse(ctx, http.StatusOK, "success", "The order has been successfully retrieved",
		&types.PaymentOrderResponse{
			ID:      paymentOrder.ID,
			Amount:  paymentOrderAmount,
			Network: paymentOrder.Edges.Token.Edges.Network.String(),
			Recipient: types.PaymentOrderRecipient{
				Institution:       paymentOrder.Edges.Recipient.Institution,
				AccountIdentifier: paymentOrder.Edges.Recipient.AccountIdentifier,
				AccountName:       paymentOrder.Edges.Recipient.AccountName,
				ProviderID:        paymentOrder.Edges.Recipient.ProviderID,
			},
			CreatedAt: paymentOrder.CreatedAt,
			UpdatedAt: paymentOrder.UpdatedAt,
			TxHash:    paymentOrder.TxHash,
			Status:    paymentOrder.Status.String(),
		})
}
