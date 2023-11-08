package controllers

import (
	"encoding/json"
	"net/http"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/routers/middleware"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/utils/test"
	"github.com/stretchr/testify/assert"
)

func setup() error {
	// Set up test data
	if _, err := test.CreateTestFiatCurrency(nil); err != nil {
		return err
	}

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

	// Set up test routers
	var ctrl Controller
	router := gin.New()
	router.Use(middleware.ScopeMiddleware)

	router.GET("currencies", ctrl.GetFiatCurrencies)

	t.Run("Currencies", func(t *testing.T) {
		t.Run("fetch supported fiat currencies", func(t *testing.T) {
			res, err := test.PerformRequest(t, "GET", "/currencies?scope=sender", nil, nil, router)
			assert.NoError(t, err)

			// Assert the response code.
			assert.Equal(t, http.StatusOK, res.Code)

			var response struct {
				Data    []types.SupportedCurrencies
				Message string
			}
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "OK", response.Message)

			// Assert /currencies response with the seeded Naira currency.
			nairaCurrency := types.SupportedCurrencies{
				Code:       "NGN",
				Name:       "Nigerian Naira",
				ShortName:  "Naira",
				Decimals:   2,
				Symbol:     "₦",
				MarketRate: decimal.NewFromFloat(950.0),
			}

			assert.Equal(t, nairaCurrency.Code, response.Data[0].Code)
			assert.Equal(t, nairaCurrency.Name, response.Data[0].Name)
			assert.Equal(t, nairaCurrency.ShortName, response.Data[0].ShortName)
			assert.Equal(t, nairaCurrency.Decimals, response.Data[0].Decimals)
			assert.Equal(t, nairaCurrency.Symbol, response.Data[0].Symbol)
			assert.True(t, response.Data[0].MarketRate.Equal(nairaCurrency.MarketRate))
		})
	})
}
