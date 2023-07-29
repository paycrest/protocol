package services

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/mattn/go-sqlite3"
	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/enttest"
	"github.com/paycrest/paycrest-protocol/types"
	"github.com/paycrest/paycrest-protocol/utils/test"
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
	backend, err := test.NewSimulatedBlockchain()
	if err != nil {
		return err
	}
	testCtx.rpcClient = backend

	// Create a test token
	token, err := test.CreateTestToken(backend, nil)
	if err != nil {
		return err
	}

	// Create receive address
	receiveAddressFactory, err := test.DeployEIP4337FactoryContract(backend)
	if err != nil {
		return err
	}
	receiveAddressService := NewReceiveAddressService()
	receiveAddress, err := receiveAddressService.CreateSmartAccount(
		context.Background(), backend, receiveAddressFactory)
	if err != nil {
		return err
	}
	testCtx.receiveAddress = receiveAddress

	// Fund receive address
	amount := big.NewInt(0)
	amount.SetString("10000000000000000000", 10)
	test.FundAddressWithTestToken(
		common.HexToAddress(token.ContractAddress),
		amount,
		common.HexToAddress(receiveAddress.Address),
	)

	// Create a payment order
	paymentOrder, err := db.Client.PaymentOrder.
		Create().
		SetAmount(decimal.NewFromBigInt(amount, 1)).
		SetAmountPaid(decimal.NewFromInt(0)).
		SetToken(token).
		SetReceiveAddress(receiveAddress).
		SetReceiveAddressText(receiveAddress.Address).
		Save(context.Background())
	if err != nil {
		return err
	}
	testCtx.paymentOrder = paymentOrder

	indexer := NewIndexerService(nil)
	testCtx.indexer = indexer

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

	// // Index ERC20 transfers for the receive address
	// done := make(chan bool)
	// err = testCtx.indexer.IndexERC20Transfer(context.Background(), testCtx.rpcClient, testCtx.receiveAddress, done)
	// assert.NoError(t, err)

	// // Fetch receiveAddress from db
	// receiveAddress, err := db.Client.ReceiveAddress.
	// 	Query().
	// 	Where(receiveaddress.AddressEQ(testCtx.receiveAddress.Address)).
	// 	Only(context.Background())
	// assert.NoError(t, err)

	// // Assert state changes after indexing
	// assert.Equal(t, receiveaddress.StatusUsed, receiveAddress.Status)
}
