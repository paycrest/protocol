package stream

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/aggregator/ent"
	db "github.com/paycrest/aggregator/storage"

	"github.com/paycrest/aggregator/ent/enttest"
	"github.com/paycrest/aggregator/utils/test"
	"github.com/stretchr/testify/assert"
)

var testCtx = struct {
	currency *ent.FiatCurrency
}{}

func setup() error {
	// Set up test data
	currency, err := test.CreateTestFiatCurrency(nil)
	if err != nil {
		return err
	}
	testCtx.currency = currency

	return nil
}

func TestIndex(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	err := setup()
	assert.NoError(t, err)

}
