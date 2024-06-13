package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/utils/test"
	"github.com/stretchr/testify/assert"
)

func setup() (*ent.FiatCurrency, error) {
	// Set up test data
	currency, err := test.CreateTestFiatCurrency(nil)
	if err != nil {
		return nil, err
	}

	return currency, nil
}

func TestIndex(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	currency, err := setup()
	assert.NoError(t, err)

	// Set up test routers
	var ctrl Controller
	router := gin.New()

	router.GET("currencies", ctrl.GetFiatCurrencies)
	router.GET("aggregator-key", ctrl.GetAggregatorPublicKey)
	router.GET("institutions/:currency_code", ctrl.GetInstitutionsByCurrency)

	t.Run("GetInstitutions By Currency", func(t *testing.T) {

		res, err := test.PerformRequest(t, "GET", fmt.Sprintf("/institutions/%s", currency.Code), nil, nil, router)
		assert.NoError(t, err)

		type Response struct {
			Status  string                        `json:"status"`
			Message string                        `json:"message"`
			Data    []types.SupportedInstitutions `json:"data"`
		}

		var response Response
		// Assert the response body
		assert.Equal(t, http.StatusOK, res.Code)

		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "OK", response.Message)
		assert.Equal(t, 2, len(response.Data), "SupportedInstitutions should be two")
	})

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

	t.Run("Get Aggregator Public key", func(t *testing.T) {
		t.Run("fetch Aggregator Public key", func(t *testing.T) {
			res, err := test.PerformRequest(t, "GET", "/aggregator-key", nil, nil, router)
			assert.NoError(t, err)

			// Assert the response code.
			assert.Equal(t, http.StatusOK, res.Code)

			var response struct {
				Data    map[string]interface{}
				Message string
			}
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "OK", response.Message)

			assert.Equal(t, response.Data["aggregatorPublicKey"], config.CryptoConfig().AggregatorPublicKey)
		})
	})
}
