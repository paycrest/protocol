package services

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/ent/receiveaddress"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var testCtx = struct {
	rpcClient      types.RPCClient
	indexer        *IndexerService
	receiveAddress *ent.ReceiveAddress
	paymentOrder   *ent.PaymentOrder
}{}

func setup() error {
	// Set up test blockchain client
	client, err := test.SetUpTestBlockchain()
	if err != nil {
		return err
	}
	testCtx.rpcClient = client

	// Create a test token
	token, err := test.CreateERC20Token(
		client,
		map[string]interface{}{})
	if err != nil {
		return err
	}

	receiveAddress, err := test.CreateSmartAccount(
		context.Background(), client)
	if err != nil {
		return err
	}
	testCtx.receiveAddress = receiveAddress

	// Create a test api key
	user, err := test.CreateTestUser(nil)
	if err != nil {
		return err
	}

	senderProfile, err := test.CreateTestSenderProfile(map[string]interface{}{
		"user_id": user.ID,
	})
	if err != nil {
		return err
	}

	apiKeyService := NewAPIKeyService()
	_, _, err = apiKeyService.GenerateAPIKey(
		context.Background(),
		nil,
		senderProfile,
		nil,
	)
	if err != nil {
		return err
	}

	// Create a payment order
	amount := decimal.NewFromFloat(29.93)
	protocolFee := amount.Mul(decimal.NewFromFloat(0.001)) // 0.1% protocol fee

	paymentOrder, err := db.Client.PaymentOrder.
		Create().
		SetSenderProfile(senderProfile).
		SetAmount(amount).
		SetAmountPaid(decimal.NewFromInt(0)).
		SetAmountReturned(decimal.NewFromInt(0)).
		SetSenderFee(decimal.NewFromInt(0)).
		SetNetworkFee(token.Edges.Network.Fee).
		SetProtocolFee(protocolFee). // 0.1% protocol fee
		SetPercentSettled(decimal.NewFromInt(0)).
		SetRate(decimal.NewFromInt(750)).
		SetToken(token).
		SetReceiveAddress(receiveAddress).
		SetReceiveAddressText(receiveAddress.Address).
		SetFeePerTokenUnit(senderProfile.FeePerTokenUnit).
		SetFeeAddress(senderProfile.FeeAddress).
		Save(context.Background())
	if err != nil {
		return err
	}
	testCtx.paymentOrder = paymentOrder

	// Fund receive address
	amountWithFees := amount.Add(paymentOrder.ProtocolFee).Add(paymentOrder.NetworkFee).Add(paymentOrder.SenderFee)
	err = test.FundAddressWithERC20Token(
		client,
		common.HexToAddress(token.ContractAddress),
		utils.ToSubunit(amountWithFees, token.Decimals),
		common.HexToAddress(receiveAddress.Address),
	)
	if err != nil {
		return err
	}

	// Create a mock instance of the OrderService
	mockOrderService := &test.MockOrderService{}

	indexer := NewIndexerService(mockOrderService)
	testCtx.indexer = indexer.(*IndexerService)

	return nil
}

func TestIndexer(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	err := setup()
	assert.NoError(t, err)

	// Index ERC20 transfers for the receive address
	_ = testCtx.indexer.IndexERC20Transfer(context.Background(), testCtx.rpcClient, testCtx.receiveAddress)
	assert.NoError(t, err)

	time.Sleep(30 * time.Second)

	// Fetch receiveAddress from db
	receiveAddress, err := db.Client.ReceiveAddress.
		Query().
		Where(receiveaddress.AddressEQ(testCtx.receiveAddress.Address)).
		Only(context.Background())
	assert.NoError(t, err)

	// Assert state changes after indexing
	assert.Equal(t, receiveaddress.StatusUsed, receiveAddress.Status)
}
