package sender

import (
	"net/http"
	"time"

	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	svc "github.com/paycrest/paycrest-protocol/services"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// Controller is a controller type for sender endpoints
type Controller struct{}

// CreatePaymentOrder controller creates a payment order
func (ctrl *Controller) CreatePaymentOrder(ctx *gin.Context) {
	var payload svc.NewPaymentOrderPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Generate receive address
	receiveAddressService := svc.NewReceiveAddressService(db.Client)

	receiveAddress, err := receiveAddressService.GenerateAndSaveAddress(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", err.Error())
		return
	}

	// Start a go routine to index the receive address
	done := make(chan bool)
	indexerService := svc.NewIndexerService(db.Client)
	go func(ctx *gin.Context, receiveAddress *ent.ReceiveAddress, done chan bool) {
		for {
			select {
			case <-done:
				return
			default:
				time.Sleep(2 * time.Minute) // add 2 minutes delay between each indexing operation

				err = indexerService.IndexERC20Transfer(ctx, receiveAddress, done)
				if err != nil {
					logger.Errorf("error: %v", err)
					return
				}
			}
		}
	}(ctx, receiveAddress, done)

	// Create payment order and recipient in a transaction
	tx, err := db.Client.Tx(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", err.Error())
		return
	}

	// Create payment order recipient
	recipient, err := tx.PaymentOrderRecipient.
		Create().
		SetInstitution(payload.Recipient.Institution).
		SetAccountIdentifier(payload.Recipient.AccountIdentifier).
		SetAccountName(payload.Recipient.AccountName).
		SetProviderID(payload.Recipient.ProviderID).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to initiate payment order", err.Error())
		_ = tx.Rollback()
		return
	}

	// Create payment order
	paymentOrder, err := tx.PaymentOrder.
		Create().
		SetAmount(payload.Amount).
		SetAmountPaid(decimal.NewFromInt(0)).
		SetNetwork(paymentorder.Network(payload.Network)).
		SetReceiveAddress(receiveAddress).
		SetReceiveAddressText(receiveAddress.Address).
		SetRecipient(recipient).
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

	u.APIResponse(ctx, http.StatusCreated, "success", "Payment order initiated successfully",
		&svc.ReceiveAddressResponse{
			ID:             paymentOrder.ID,
			Amount:         paymentOrder.Amount,
			Network:        paymentOrder.Network.String(),
			ReceiveAddress: paymentOrder.ReceiveAddressText,
		})
}

// GetPaymentOrderByID controller fetches a payment order by ID
func (ctrl *Controller) GetPaymentOrderByID(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// DeletePaymentOrder controller deletes a payment order
func (ctrl *Controller) DeletePaymentOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}
