package utils

import (
	"fmt"

	"strings"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/shopspring/decimal"
)

var (
	BitgetAPIURL  = "https://api.bitget.com"
	BinanceAPIURL = "https://api.binance.com"
	QuidaxAPIURL  = "https://www.quidax.com/api/v1"
)

// fetchExternalRate fetches the external rate for a fiat currency
func FetchExternalRate(currency string) (decimal.Decimal, error) {
	currency = strings.ToUpper(currency)
	supportedCurrencies := []string{"KES", "NGN", "GHS", "TZS", "UGX", "XOF"}
	isSupported := false
	for _, supported := range supportedCurrencies {
		if currency == supported {
			isSupported = true
			break
		}
	}
	if !isSupported {
		return decimal.Zero, fmt.Errorf("ComputeMarketRate: currency not supported")
	}

	var prices []decimal.Decimal

	// Fetch rates based on currency
	if currency == "NGN" {
		quidaxRate, err := FetchQuidaxRate(currency)
		if err == nil {
			prices = append(prices, quidaxRate)
		}
	} else {
		binanceRate, err := FetchBinanceRate(currency)
		if err == nil {
			prices = append(prices, binanceRate)
		}
	}

	// Fetch Bitget rate for all supported currencies
	bitgetRate, err := FetchBitgetRate(currency)
	if err == nil {
		prices = append(prices, bitgetRate)
	}

	if len(prices) == 0 {
		return decimal.Zero, fmt.Errorf("ComputeMarketRate: no valid rates found")
	}

	// Return the median price
	return Median(prices), nil
}

// FetchQuidaxRate fetches the USDT exchange rate from Quidax (NGN only)
func FetchQuidaxRate(currency string) (decimal.Decimal, error) {
	url := fmt.Sprintf("/api/v1/markets/tickers/usdt%s", strings.ToLower(currency))

	res, err := fastshot.NewClient(QuidaxAPIURL).
		Config().SetTimeout(30*time.Second).
		Build().GET(url).
		Retry().Set(3, 5*time.Second).
		Send()
	if err != nil {
		return decimal.Zero, fmt.Errorf("FetchQuidaxRate: %w", err)
	}

	data, err := ParseJSONResponse(res.RawResponse)
	if err != nil {
		return decimal.Zero, fmt.Errorf("FetchQuidaxRate: %w", err)
	}

	price, err := decimal.NewFromString(data["data"].(map[string]interface{})["ticker"].(map[string]interface{})["buy"].(string))
	if err != nil {
		return decimal.Zero, fmt.Errorf("FetchQuidaxRate: %w", err)
	}

	return price, nil
}

// FetchBinanceRate fetches the median USDT exchange rate from Binance P2P
func FetchBinanceRate(currency string) (decimal.Decimal, error) {

	res, err := fastshot.NewClient(BinanceAPIURL).
		Config().SetTimeout(30*time.Second).
		Header().Add("Content-Type", "application/json").
		Build().POST("/bapi/c2c/v2/friendly/c2c/adv/search").
		Retry().Set(3, 5*time.Second).
		Body().AsJSON(map[string]interface{}{
		"asset":     "USDT",
		"fiat":      currency,
		"tradeType": "SELL",
		"page":      1,
		"rows":      20,
	}).
		Send()
	if err != nil {
		return decimal.Zero, fmt.Errorf("FetchBinanceRate: %w", err)
	}

	resData, err := ParseJSONResponse(res.RawResponse)
	if err != nil {
		return decimal.Zero, fmt.Errorf("FetchBinanceRate: %w", err)
	}

	data, ok := resData["data"].([]interface{})
	if !ok || len(data) == 0 {
		return decimal.Zero, fmt.Errorf("FetchBinanceRate: no data in response")
	}

	var prices []decimal.Decimal
	for _, item := range data {
		adv, ok := item.(map[string]interface{})["adv"].(map[string]interface{})
		if !ok {
			continue
		}

		price, err := decimal.NewFromString(adv["price"].(string))
		if err != nil {
			continue
		}

		prices = append(prices, price)
	}

	if len(prices) == 0 {
		return decimal.Zero, fmt.Errorf("FetchBinanceRate: no valid prices found")
	}

	return Median(prices), nil
}

// FetchBitgetRate fetches the median USDT exchange rate from Bitget P2P
func FetchBitgetRate(currency string) (decimal.Decimal, error) {

	res, err := fastshot.NewClient(BitgetAPIURL).
		Config().SetTimeout(30*time.Second).
		Header().Add("Content-Type", "application/json").
		Build().POST("/api/v2/p2p/adv/search").
		Retry().Set(3, 5*time.Second).
		Body().AsJSON(map[string]interface{}{
		"tokenId": "USDT",
		"fiat":    currency,
		"side":    "sell",
		"page":    1,
		"size":    20,
	}).
		Send()
	if err != nil {
		return decimal.Zero, fmt.Errorf("FetchBitgetRate: %w", err)
	}

	resData, err := ParseJSONResponse(res.RawResponse)
	if err != nil {
		return decimal.Zero, fmt.Errorf("FetchBitgetRate: %w", err)
	}

	data, ok := resData["data"].([]interface{})
	if !ok || len(data) == 0 {
		return decimal.Zero, fmt.Errorf("FetchBitgetRate: no data in response")
	}

	var prices []decimal.Decimal
	for _, item := range data {
		adv, ok := item.(map[string]interface{})["price"].(string)
		if !ok {
			continue
		}

		price, err := decimal.NewFromString(adv)
		if err != nil {
			continue
		}

		prices = append(prices, price)
	}

	if len(prices) == 0 {
		return decimal.Zero, fmt.Errorf("FetchBitgetRate: no valid prices found")
	}

	return Median(prices), nil
}
