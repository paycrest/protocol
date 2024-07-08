package order

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/enttest"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils/test"
	"github.com/stretchr/testify/assert"
)

var testCxtEVM = struct {
	blockchainClient types.RPCClient
	user             *ent.User
	paymentOrder     *ent.PaymentOrder
	currency         *ent.FiatCurrency
	client           types.RPCClient
	token            *ent.Token
}{}

func setupEVM() error {

	backend, err := test.SetUpTestBlockchain()
	if err != nil {
		return err
	}

	testCxtEVM.blockchainClient = backend

	token, err := test.CreateERC20Token(backend, map[string]interface{}{})
	if err != nil {
		return err
	}

	testCxtEVM.token = token

	user, err := test.CreateTestUser(map[string]interface{}{
		"scope": "provider",
		"email": "providerjohndoe@test.com",
	})
	if err != nil {
		return err
	}

	testCxtEVM.user = user

	currency, err := test.CreateTestFiatCurrency(map[string]interface{}{
		"code":        "KES",
		"short_name":  "Shilling",
		"decimals":    2,
		"symbol":      "KSh",
		"name":        "Kenyan Shilling",
		"market_rate": 550.0,
	})
	if err != nil {
		return err
	}

	testCxtEVM.currency = currency

	sender, err := test.CreateTestSenderProfile(map[string]interface{}{
		"user_id":     user.ID,
		"currency_id": currency.ID,
	})
	if err != nil {
		return err
	}

	paymentOrder, err := test.CreateTestPaymentOrder(backend, token, map[string]interface{}{
		"sender": sender,
	})

	if err != nil {
		return err
	}
	testCxtEVM.paymentOrder = paymentOrder

	return nil

}

func TestEVMOrders(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	err := setupEVM()
	assert.NoError(t, err)

	orderservice := NewOrderEVM()

	t.Run("createOrder", func(t *testing.T) {
		err = orderservice.CreateOrder(context.Background(), testCxtEVM.paymentOrder.ID)
		assert.NoError(t, err)
	})
}
