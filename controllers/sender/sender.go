package sender

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/paycrest/protocol/ent"
	db "github.com/paycrest/protocol/storage"

	"github.com/paycrest/protocol/ent/network"
	"github.com/paycrest/protocol/ent/paymentorder"
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

	// Get sender profile from the context
	senderCtx, ok := ctx.Get("sender")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	sender := senderCtx.(*ent.SenderProfile)

	// Get token from DB
	token, err := db.Client.Token.
		Query().
		Where(
			token.SymbolEQ(payload.Token),
			token.HasNetworkWith(network.IdentifierEQ(payload.Network)),
		).
		WithNetwork().
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Provided token is not supported", nil)
		return
	}

	// Generate receive address
	receiveAddress, err := ctrl.receiveAddressService.CreateSmartAccount(ctx, nil, nil)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", nil)
		return
	}

	// Create payment order and recipient in a transaction
	tx, err := db.Client.Tx(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", nil)
		return
	}

	// Create payment order
	paymentOrder, err := tx.PaymentOrder.
		Create().
		SetSenderProfile(sender).
		SetAmount(payload.Amount).
		SetAmountPaid(decimal.NewFromInt(0)).
		SetToken(token).
		SetLabel(payload.Label).
		SetRate(payload.Rate).
		SetReceiveAddress(receiveAddress).
		SetReceiveAddressText(receiveAddress.Address).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", nil)
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
		SetMemo(payload.Recipient.Memo).
		SetPaymentOrder(paymentOrder).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", nil)
		_ = tx.Rollback()
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", nil)
		return
	}

	// Start a background process to index token transfers to the receive address
	go ctrl.indexerService.RunIndexERC20Transfer(ctx, receiveAddress)

	paymentOrderAmount, _ := paymentOrder.Amount.Float64()

	u.APIResponse(ctx, http.StatusCreated, "success", "Payment order initiated successfully",
		&types.ReceiveAddressResponse{
			ID:             paymentOrder.ID,
			Amount:         paymentOrderAmount,
			Network:        token.Edges.Network.Identifier,
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
			"Invalid order ID", nil)
		return
	}

	// Get sender profile from the context
	senderCtx, ok := ctx.Get("sender")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	sender := senderCtx.(*ent.SenderProfile)

	// Fetch payment order from the database
	paymentOrder, err := db.Client.PaymentOrder.
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
		ID:      paymentOrder.ID,
		Amount:  paymentOrder.Amount,
		Rate:    paymentOrder.Rate,
		Network: paymentOrder.Edges.Token.Edges.Network.Identifier,
		Recipient: types.PaymentOrderRecipient{
			Institution:       paymentOrder.Edges.Recipient.Institution,
			AccountIdentifier: paymentOrder.Edges.Recipient.AccountIdentifier,
			AccountName:       paymentOrder.Edges.Recipient.AccountName,
			ProviderID:        paymentOrder.Edges.Recipient.ProviderID,
			Memo:              paymentOrder.Edges.Recipient.Memo,
		},
		Label:     paymentOrder.Label,
		CreatedAt: paymentOrder.CreatedAt,
		UpdatedAt: paymentOrder.UpdatedAt,
		TxHash:    paymentOrder.TxHash,
		Status:    paymentOrder.Status,
	})
}
