package services

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/ent/receiveaddress"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
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
	backend, err := test.NewSimulatedBlockchain()
	if err != nil {
		return err
	}
	testCtx.rpcClient = backend

	// TODO: replace this implementation with the simulated blockchain
	// We're currently using this tx for testing:
	// https://etherscan.io/tx/0xd5fb7204066619f207e8d2c87493cdd1edbdd9ec21ab647f95d45faa2ed540cb

	// Create a test token
	token, err := test.CreateTestToken(
		backend,
		map[string]interface{}{
			"symbol":           "USDT",
			"contract_address": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
			"decimals":         6,
			"networkRPC":       "https://mainnet.infura.io/v3/4818dbcee84d4651a832894818bd4534",
		})
	if err != nil {
		return err
	}

	// Create receive address
	// receiveAddressFactory, err := test.DeployEIP4337FactoryContract(backend)
	// if err != nil {
	// 	return err
	// }
	// receiveAddressService := NewReceiveAddressService()
	// receiveAddress, err := receiveAddressService.CreateSmartAccount(
	// 	context.Background(), nil, nil)
	// if err != nil {
	// 	return err
	// }

	// Save address in db
	receiveAddress, err := db.Client.ReceiveAddress.
		Create().
		SetAddress("0xF6F6407410235202CA5Bfa68286a3bBe01F8E5E0").
		SetSalt([]byte("random salt")).
		SetStatus(receiveaddress.StatusUnused).
		SetLastIndexedBlock(17800411). // our target test block with usdt transfer is 17800412
		Save(context.Background())
	if err != nil {
		return err
	}
	testCtx.receiveAddress = receiveAddress

	// Fund receive address
	amount := decimal.NewFromInt(2990)
	// err = test.FundAddressWithTestToken(
	// 	backend,
	// 	common.HexToAddress(token.ContractAddress),
	// 	amount,
	// 	common.HexToAddress(receiveAddress.Address),
	// )
	// if err != nil {
	// 	return err
	// }

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
	paymentOrder, err := db.Client.PaymentOrder.
		Create().
		SetSenderProfile(senderProfile).
		SetAmount(amount).
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

	// Index ERC20 transfers for the receive address
	_, err = testCtx.indexer.IndexERC20Transfer(context.Background(), nil, testCtx.receiveAddress)
	assert.NoError(t, err)

	// Fetch receiveAddress from db
	receiveAddress, err := db.Client.ReceiveAddress.
		Query().
		Where(receiveaddress.AddressEQ(testCtx.receiveAddress.Address)).
		Only(context.Background())
	assert.NoError(t, err)

	// Assert state changes after indexing
	assert.Equal(t, receiveaddress.StatusUsed, receiveAddress.Status)

}
